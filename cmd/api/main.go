package main

import (
	"fmt"
	"log"
	"meme-trader/internal/api"
	"meme-trader/internal/config"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create and initialize the application
	app, err := api.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on %s", addr)
	if err := app.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
