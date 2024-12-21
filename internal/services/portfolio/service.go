package portfolio

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/postgres"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
)

type Service struct {
	db    *postgres.PortfolioRepository
	cache redis.Cache
}

func NewService(db *postgres.PortfolioRepository, cache redis.Cache) *Service {
	return &Service{db: db, cache: cache}
}

func (s *Service) GetHoldings(ctx context.Context, userID string) ([]models.PortfolioHolding, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	holdings, err := s.db.GetHoldings(ctx, uid)
	if err != nil {
		return nil, err
	}

	return holdings, nil
}

func (s *Service) GetHistory(ctx context.Context, userID string, timeframe string) ([]models.PortfolioValue, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	history, err := s.db.GetHistory(ctx, uid, timeframe)
	if err != nil {
		return nil, err
	}

	return history, nil
}
