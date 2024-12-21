package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nicolas-martin/memecoin-trading/internal/api/handlers"
	"github.com/nicolas-martin/memecoin-trading/internal/config"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/postgres"
	redisRepo "github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
	"github.com/nicolas-martin/memecoin-trading/internal/services/coin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/leaderboard"
	"github.com/nicolas-martin/memecoin-trading/internal/services/payment"
	"github.com/nicolas-martin/memecoin-trading/internal/services/portfolio"
	"github.com/nicolas-martin/memecoin-trading/internal/services/support"
	"github.com/nicolas-martin/memecoin-trading/pkg/dexscreens"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	router *gin.Engine
	db     *gorm.DB
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := postgres.NewDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &App{
		config: cfg,
		router: gin.Default(),
		db:     db,
	}, nil
}

func (a *App) Start() error {
	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", a.config.Redis.Host, a.config.Redis.Port),
		Password: a.config.Redis.Password,
		DB:       a.config.Redis.DB,
	})

	cache := redisRepo.NewRedisCache(redisClient)

	// Initialize repositories
	leaderboardRepo := postgres.NewLeaderboardRepository(a.db)
	portfolioRepo := postgres.NewPortfolioRepository(a.db)
	supportRepo := postgres.NewSupportRepository(a.db)
	paymentRepo := postgres.NewPaymentRepository(a.db)

	// Initialize DexScreens client
	dexScreens := dexscreens.NewClient(
		a.config.DexScreens.ApiURL,
		a.config.DexScreens.ApiKey,
	)

	// Initialize services
	coinService := coin.NewService(cache, dexScreens)
	leaderboardService := leaderboard.NewService(leaderboardRepo, cache)
	portfolioService := portfolio.NewService(portfolioRepo, cache)
	supportService := support.NewService(supportRepo)
	paymentService := payment.NewService(paymentRepo, cache)

	// Initialize handlers
	handler := handlers.NewHandler(
		leaderboardService,
		portfolioService,
		supportService,
		coinService,
		paymentService,
	)

	// Setup routes
	setupRoutes(a.router, handler)

	// Start server
	return a.router.Run(fmt.Sprintf(":%s", a.config.App.Port))
}
