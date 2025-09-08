package migration

import (
	"fmt"
	"go-microservices/customer-service/models"
	"go-microservices/pkg/config"
	"log"
)

func AutoMigrate() {
    err := config.DB.AutoMigrate(&models.Customer{})
    if err != nil {
        log.Fatalf("failed to migrate: %v", err)
    }	
	fmt.Println("Gorm DB Setup successfull !")
}