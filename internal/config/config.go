package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	SolanaEndpoint string
	DatabaseURL    string
}

func NewConfig() *Config {
	// Load .env.local first, then fall back to .env
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: No .env or .env.local file found")
		}
	}

	// Build database URL if not provided directly
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			getEnvOrDefault("DB_USER", ""),
			getEnvOrDefault("DB_PASSWORD", ""),
			getEnvOrDefault("DB_HOST", "localhost"),
			getEnvOrDefault("DB_PORT", "5432"),
			getEnvOrDefault("DB_NAME", "memetrader"),
		)
	}

	return &Config{
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		SolanaEndpoint: getEnvOrDefault("SOLANA_ENDPOINT", "https://api.devnet.solana.com"),
		DatabaseURL:    dbURL,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate checks if all required environment variables are set
func (c *Config) Validate() error {
	if c.DatabaseURL == "postgres://:@localhost:5432/memetrader?sslmode=disable" {
		return fmt.Errorf("database credentials not provided. Please set DATABASE_URL or individual DB_* environment variables")
	}
	return nil
}
