package main

import (
	"log"
	"net/http"
	"os"

	"go-microservices/gateway/routes"
	initializers "go-microservices/pkg/initializer"

	"github.com/gin-gonic/gin"
)

func init() {
    initializers.LoadEnvvariables()
}

func main() {

    r := gin.Default()

    api := r.Group("/api/v1")
    {
        routes.CustomerRoutes(api) // attach customer routes
    }

    port := os.Getenv("GATEWAY_PORT")
    log.Printf("🚀 API Gateway running on  %s", port)
    if err := r.Run(port); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
