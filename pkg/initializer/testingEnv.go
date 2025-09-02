package initializers

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func TestLoadEnvvariables() {
	// Absolute path to project root .env
	envPath, _ := filepath.Abs(filepath.Join("..", "..", ".env"))

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s", envPath)
	}
}
