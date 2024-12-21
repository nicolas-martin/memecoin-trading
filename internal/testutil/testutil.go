package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB creates a test database connection
func TestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=memecoin_trading_test port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

// TestRedis creates a test Redis connection
func TestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use different DB for tests
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("Failed to connect to test Redis: %v", err)
	}

	return client
}

// CleanupDB cleans up the test database
func CleanupDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("Failed to get underlying sql.DB: %v", err)
		return
	}
	sqlDB.Close()
}

// CleanupRedis cleans up the test Redis database
func CleanupRedis(t *testing.T, client *redis.Client) {
	if err := client.FlushDB(context.Background()).Err(); err != nil {
		t.Errorf("Failed to flush Redis DB: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Errorf("Failed to close Redis connection: %v", err)
	}
}
