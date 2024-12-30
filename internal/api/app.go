package api

import (
	"context"
	"log"
	"meme-trader/internal/api/handlers"
	"meme-trader/internal/config"
	"meme-trader/internal/repository/postgres"
	"meme-trader/internal/services/memecoin"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Router  *mux.Router
	Service *memecoin.Service
}

func NewApp(cfg *config.Config) (*App, error) {
	// Initialize database
	db, err := postgres.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Create logger
	logger := log.New(os.Stdout, "memecoin-service: ", log.LstdFlags)

	// Initialize service
	service := memecoin.NewService(db, logger)

	// Create router
	router := mux.NewRouter()

	// Initialize handlers
	memeHandler := handlers.NewMemeHandler(service)

	// Register routes
	router.HandleFunc("/api/v1/memecoins", memeHandler.GetTopMemeCoins).Methods("GET")
	router.HandleFunc("/api/v1/memecoins/{id}", memeHandler.GetMemeCoinDetail).Methods("GET")
	router.HandleFunc("/api/v1/memecoins/update", func(w http.ResponseWriter, r *http.Request) {
		if err := service.FetchAndUpdateMemeCoins(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	return &App{
		Router:  router,
		Service: service,
	}, nil
}

func (a *App) Run(addr string) error {
	// Initial fetch of meme coins
	if err := a.Service.FetchAndUpdateMemeCoins(context.Background()); err != nil {
		return err
	}

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",  // Expo web
			"http://localhost:19006", // Expo web alternative port
			"http://localhost:3000",  // Common React port
			"exp://localhost:8081",   // Expo development
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		Debug:          true, // Enable debugging for development
	})

	// Wrap router with CORS middleware
	handler := c.Handler(a.Router)
	return http.ListenAndServe(addr, handler)
}
