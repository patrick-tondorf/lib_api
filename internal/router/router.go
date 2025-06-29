package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/patrick-tondorf/lib_api/docs"
	"github.com/patrick-tondorf/lib_api/internal/config"
	"github.com/patrick-tondorf/lib_api/internal/handler"
	auth "github.com/patrick-tondorf/lib_api/internal/middleware"
	"github.com/patrick-tondorf/lib_api/internal/repository"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *pgx.Conn) *gin.Engine {
	r := gin.New()

	// Middlewares básicos
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configuração do Swagger
	docs.SwaggerInfo.Title = "Library API"
	docs.SwaggerInfo.Description = "Virtual Library API with Supabase and JWT Authentication"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Rotas do Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Inicializa repositórios e handlers
	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)
	authorRepo := repository.NewAuthorRepository(db)
	authorHandler := handler.NewAuthorHandler(authorRepo)
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Secret key for JWT - agora com fallback para config.GetSecretKey()
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = config.GetSecretKey() // Adicione esta função no seu pacote config
		if secret == "" {
			panic("SECRET_KEY not configured in environment variables or config")
		}
	}

	// Rotas públicas
	public := r.Group("/api")
	{
		// User routes
		public.POST("/users", userHandler.CreateUser)
		public.POST(("/auth/login"), userHandler.AuthenticateUser)
		// Rotas públicas adicionais (se houver)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	// Rotas protegidas
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware(secret))
	{

		//user routes
		protected.GET("/users/:email", userHandler.GetUserByEmail)
		// Book routes
		protected.POST("/books", bookHandler.CreateBook)
		protected.GET("/books", bookHandler.GetBooks)
		//protected.GET("/books/:id", bookHandler.GetBookByID)
		protected.PUT("/books/:id", bookHandler.UpdateBook)
		//protected.DELETE("/books/:id", bookHandler.DeleteBook)

		// Author routes
		protected.POST("/authors", authorHandler.CreateAuthor)
		protected.GET("/authors", authorHandler.GetAuthors)
		protected.GET("/authors/:id", authorHandler.GetAuthorByID)
		//protected.PUT("/authors/:id", authorHandler.UpdateAuthor)
		//protected.DELETE("/authors/:id", authorHandler.DeleteAuthor)

		// Rotas protegidas adicionais do usuário
		//protected.GET("/users/me", userHandler.GetCurrentUser)
		//protected.PUT("/users/me", userHandler.UpdateCurrentUser)
	}

	return r
}
