package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	SolanaEndpoint string
}

func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))

	return &Config{
		DBHost:         getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:         dbPort,
		DBUser:         getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:     getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:         getEnvOrDefault("DB_NAME", "memetrader"),
		ServerPort:     getEnvOrDefault("SERVER_PORT", "8080"),
		SolanaEndpoint: getEnvOrDefault("SOLANA_ENDPOINT", "https://api.devnet.solana.com"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
