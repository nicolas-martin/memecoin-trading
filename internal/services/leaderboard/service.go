package leaderboard

import (
	"context"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/postgres"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
)

type Service struct {
	db    *postgres.LeaderboardRepository
	cache redis.Cache
}

func NewService(db *postgres.LeaderboardRepository, cache redis.Cache) *Service {
	return &Service{db: db, cache: cache}
}

func (s *Service) GetTopTraders(ctx context.Context, duration time.Duration) ([]models.LeaderboardEntry, error) {
	timeframe := duration.String()

	// Try to get from cache first
	if entries, err := s.cache.GetLeaderboard(ctx, timeframe); err == nil && entries != nil {
		return entries, nil
	}

	// If not in cache, get from database
	entries, err := s.db.GetTopTraders(ctx, duration, 100)
	if err != nil {
		return nil, err
	}

	// Cache the results
	if err := s.cache.SetLeaderboard(ctx, timeframe, entries); err != nil {
		// Log error but don't fail the request
		// logger.Error("Failed to cache leaderboard", "error", err)
	}

	return entries, nil
}
