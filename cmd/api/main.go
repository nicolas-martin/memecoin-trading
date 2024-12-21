package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nicolas-martin/memecoin-trading/internal/api"
	"github.com/nicolas-martin/memecoin-trading/internal/config"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize app
	app, err := api.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// Start the server
	if err := app.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
