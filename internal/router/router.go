package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/patrick-tondorf/lib_api/docs"
	"github.com/patrick-tondorf/lib_api/internal/handler"
	"github.com/patrick-tondorf/lib_api/internal/repository"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *pgx.Conn) *gin.Engine {
	r := gin.New()

	// Adicione manualmente os middlewares que precisar
	r.Use(gin.Logger())   // Se quiser o logger
	r.Use(gin.Recovery()) // Se quiser o recovery

	//Config Swagger
	// Configuração do Swagger
	docs.SwaggerInfo.Title = "Library API"
	docs.SwaggerInfo.Description = "Virtual Library API with Supabase"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	// Rotas do Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//Inicializa repositório e handlers
	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	//Rotas
	api := r.Group("/api")
	{
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books", bookHandler.GetBooks)
	}
	return r
}
