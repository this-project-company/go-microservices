package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

// InitRedis initializes the Redis client
func InitRedis() {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	addr :=  os.Getenv("REDIS_ADDR")
	pass := os.Getenv("REDIS_PASSWORD")
	
	Rdb = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", addr),
		Password: getEnv("REDIS_PASSWORD", pass),
		DB:       db,
	})

	// Test connection
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis")
}

// Set stores a key with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	return Rdb.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a key
func Get(key string) (string, error) {
	return Rdb.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(key string) error {
	return Rdb.Del(ctx, key).Err()
}

// getEnv helper
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
