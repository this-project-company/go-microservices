package config

import (
	"os"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GormDB() {
	dsn := os.Getenv("GORM_DB_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	
	// Get the generic DB object from gorm.DB to configure pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// ✅ Connection pool settings
	sqlDB.SetMaxOpenConns(25)                  // max number of open connections
	sqlDB.SetMaxIdleConns(5)                   // max number of idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute)  // max lifetime of a connection
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)  // max idle time of a connection

	DB = db
	fmt.Println("Database connected with pool configured in Gorm!")

}


func CloseGorm() {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err == nil {
            sqlDB.Close()
            fmt.Println("Database connection closed!")
        }
    }
}