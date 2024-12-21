package redis

import (
	"context"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
)

type Cache interface {
	// Coin methods
	GetTopCoins(ctx context.Context, limit int) ([]models.Coin, error)
	SetTopCoins(ctx context.Context, coins []models.Coin) error
	GetCoinByID(ctx context.Context, id string) (*models.Coin, error)
	SetCoin(ctx context.Context, coin *models.Coin) error
	InvalidateCoinCache(ctx context.Context, id string) error

	// User methods
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	SetUser(ctx context.Context, user *models.User) error
	InvalidateUserCache(ctx context.Context, id string) error

	// Transaction stats methods
	GetUserStats(ctx context.Context, userID string) (map[string]float64, error)
	SetUserStats(ctx context.Context, userID string, stats map[string]float64) error
	InvalidateUserStats(ctx context.Context, userID string) error

	// Leaderboard methods
	GetLeaderboard(ctx context.Context, timeframe string) ([]models.LeaderboardEntry, error)
	SetLeaderboard(ctx context.Context, timeframe string, entries []models.LeaderboardEntry) error
	InvalidateLeaderboard(ctx context.Context, timeframe string) error
}
