package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrick-tondorf/lib_api/internal/config"
	"github.com/patrick-tondorf/lib_api/internal/domain"
	"github.com/patrick-tondorf/lib_api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.UserRepository
}

// LoginResponse defines the structure of a successful login response
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body domain.User true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validação básica de email
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Gera o hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}

	// Prepara o usuário para o banco
	dbUser := domain.User{
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
		// ID será gerado pelo Supabase
	}

	if err := h.repo.CreateUser(c.Request.Context(), dbUser); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func isValidEmail(email string) bool {
	// Implemente uma validação simples ou use regex
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// User godoc
// @Summary Authenticate a user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.Credentials true "User credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	log.Println("Iniciando autenticação do usuário")

	var credentials domain.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("Erro ao decodificar JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Printf("Tentativa de login para o email: %s\n", credentials.Email)

	user, err := h.repo.GetUserByEmail(c.Request.Context(), credentials.Email)
	if err != nil {
		log.Printf("Erro ao buscar usuário por email: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	log.Println("Usuário encontrado no banco de dados")

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Println("Senha incorreta para o usuário")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}
		log.Printf("Erro ao comparar senhas: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": "password comparison failed",
		})
		return
	}

	log.Println("Credenciais validadas com sucesso")

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = config.GetSecretKey()
		log.Println("Usando secret key do config")
	} else {
		log.Println("Usando secret key da variável de ambiente")
	}

	if secret == "" {
		log.Println("Nenhuma secret key configurada")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": "JWT secret not configured",
		})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Erro ao gerar token JWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"details": "JWT signing failed",
		})
		return
	}

	log.Println("Token JWT gerado com sucesso")

	response := LoginResponse{
		Token:     "Bearer " + tokenString,
		ExpiresIn: int64(time.Hour.Seconds() * 24),
		TokenType: "Bearer",
	}

	c.JSON(http.StatusOK, response)
}

// User godoc
// @Summary Get user by email
// @Description Retrieve user details by email
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param email path string true "User email" example("user@example.com")// @Success 200 {object} domain.User
// @Failure 400 {object} UserListResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{email} [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
