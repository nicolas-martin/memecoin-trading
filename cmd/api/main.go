package main

import (
	"log"
	"meme-trader/internal/api"
	"meme-trader/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	app, err := api.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := app.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
