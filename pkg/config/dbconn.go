package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

// DBconnection initializes the global pgx pool
func DBconnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := os.Getenv("DB_URL")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	// Optional pool tuning
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	DBPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect: %v\n", err)
	}

	fmt.Println("✅ Connected to Postgres with pgxpool")
}

// CloseDB closes the connection pool gracefully
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		fmt.Println("🛑 Postgres connection closed")
	}
}
