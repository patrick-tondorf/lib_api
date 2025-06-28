package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/patrick-tondorf/lib_api/internal/config"
	"github.com/patrick-tondorf/lib_api/internal/router"
)

func main() {
	// Carrega o arquivo .env
	err := godotenv.Load() // Procura por .env na raiz
	if err != nil {
		log.Fatal("Erro ao carregar .env", err)
	}

	// Config gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Acessa as variáveis
	dbURI := os.Getenv("DB_URI")
	port := os.Getenv("LOCAL_PORT")
	if port == "" {
		port = "8080" // Valor padrão se não existir no .env
	}
	log.Println("DB_URI:", dbURI)

	//Conecta ao Supabase
	db, err := config.NewSupabaseDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	//Inicia o router
	r := router.SetupRouter(db)

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Erro ao configurar proxies confiáveis:", err)
	}

	// Inicia Servidor
	log.Printf("Servidor rodando na porta %s (modo: %s)", port, gin.Mode())
	printRoutes(r)               // Exibe todas as rotas no console
	err_run := r.Run(":" + port) // localHost:8080
	if err_run != nil {
		log.Fatal("Error ao iniciar o servidor", err_run)
	}

}
func printRoutes(r *gin.Engine) {
	for _, route := range r.Routes() {
		fmt.Printf("Method: %-6s | Path: %s\n", route.Method, route.Path)
	}
}
