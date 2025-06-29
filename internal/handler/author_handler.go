package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrick-tondorf/lib_api/internal/domain"
	"github.com/patrick-tondorf/lib_api/internal/repository"
)

type AuthorHandler struct {
	Repo *repository.AuthorRepository
}

func NewAuthorHandler(repo *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{Repo: repo}
}

// CreateAuthor godoc
// @Summary Create a new author
// @Description Create a new author in the system
// @Tags authors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param author body domain.Author false "Author data"
// @Success 201 {object} domain.Author
// @Failure 400 {object} domain.Author
// @Failure 500 {object} domain.Author
// @Router /authors [post]
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var input domain.Author

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	author := domain.Author{Name: input.Name}

	if err := h.Repo.CreateAuthor(c.Request.Context(), &author); err != nil {
		log.Printf("Error creating author: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	c.JSON(http.StatusCreated, author) // Já contém todos os campos
}

// ListAll godoc
// @Summary List all authors
// @Description Get all authors from the library
// @Tags authors
// @Security BearerAuth
// @Produce json
// @Param withBooks query boolean false "Include books in response"
// @Success 200 {array} domain.Author
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /authors [get]
func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	withBooks := c.Query("withBooks") == "true"

	var authors []domain.Author
	var err error

	if withBooks {
		authors, err = h.Repo.GetAuthorsWithBooks(c.Request.Context())
	} else {
		authors, err = h.Repo.GetAuthors(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch authors",
			"details": err.Error(),
		})
		return
	}

	if len(authors) == 0 {
		c.JSON(http.StatusOK, []domain.Author{})
		return
	}

	c.JSON(http.StatusOK, authors)
}

// GetAuthorById godoc
// @Summary Get an author by ID
// @Description Retrieve an author by their ID
// @Tags authors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} domain.Author
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /authors/{id} [get]
func (h *AuthorHandler) GetAuthorByID(c *gin.Context) {
	/*id := c.Param("id")
	author, err := h.repo.GetAuthorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	c.JSON(http.StatusOK, author)*/
}
