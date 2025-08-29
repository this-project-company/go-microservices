package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvvariables() {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ No .env file found, relying on system environment variables")
    }
}
