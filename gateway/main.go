package main

import (
	"log"
	"net/http"

	"go-microservices/gateway/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    api := r.Group("/api/v1")
    {
        routes.CustomerRoutes(api) // attach customer routes
    }

    log.Println("🚀 API Gateway running on :8080")
    if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
