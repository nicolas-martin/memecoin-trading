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
	router.HandleFunc("/api/v1/memecoins", memeHandler.GetTopMemeCoins).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/memecoins/{id}", memeHandler.GetMemeCoinDetail).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/memecoins/update", func(w http.ResponseWriter, r *http.Request) {
		if err := service.FetchAndUpdateMemeCoins(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("POST", "OPTIONS")

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
			"http://localhost:8080",  // API server
			"http://localhost:8081",  // Expo web
			"http://localhost:8082",  // Expo web alternative
			"http://localhost:19000", // Expo dev tools
			"http://localhost:19001", // Expo dev tools alternative
			"http://localhost:19002", // Expo dev tools alternative
			"http://localhost:19006", // Expo web alternative
			"http://localhost:3000",  // Common React port
			"exp://localhost:8081",   // Expo development
			"http://127.0.0.1:8082",  // Local IP alternative
			"http://127.0.0.1:8081",  // Local IP alternative
			"http://127.0.0.1:19006", // Local IP alternative
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		Debug:          true, // Enable debugging for development
	})

	// Wrap router with CORS middleware
	handler := c.Handler(a.Router)

	log.Printf("Server starting on %s with CORS enabled for development", addr)
	return http.ListenAndServe(addr, handler)
}
