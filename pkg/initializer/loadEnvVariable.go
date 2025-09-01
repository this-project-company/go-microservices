package initializers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvvariables() {
    // go one level up from customer-service/
    root, _ := os.Getwd()
    envPath := filepath.Join(root, "..", ".env")

    err := godotenv.Load(envPath)
    if err != nil {
        log.Fatalf("Error loading .env file from %s", envPath)
    }
}
