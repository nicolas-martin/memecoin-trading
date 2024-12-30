package api

import (
	"fmt"
	"meme-trader/internal/api/handlers"
	"meme-trader/internal/config"
	"meme-trader/internal/repository/postgres"
	"meme-trader/internal/services/memecoin"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Router  *mux.Router
	DB      *postgres.Database
	Service *memecoin.Service
	Config  *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := postgres.NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	service := memecoin.NewService(db)
	router := mux.NewRouter()

	app := &App{
		Router:  router,
		DB:      db,
		Service: service,
		Config:  cfg,
	}

	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	memeHandler := handlers.NewMemeHandler(a.Service)

	// API routes
	api := a.Router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/memecoins", memeHandler.GetTopMemeCoins).Methods("GET")
	api.HandleFunc("/memecoins/{id}", memeHandler.GetMemeCoinDetail).Methods("GET")
	api.HandleFunc("/memecoins/update", memeHandler.UpdateMemeCoins).Methods("POST")

	// Start background job to update meme coins
	go a.startUpdateJob()
}

func (a *App) startUpdateJob() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		if err := a.Service.FetchAndUpdateMemeCoins(); err != nil {
			fmt.Printf("Error updating meme coins: %v\n", err)
		}
	}
}

func (a *App) Run() error {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(a.Router)
	return http.ListenAndServe(":"+a.Config.ServerPort, handler)
}
