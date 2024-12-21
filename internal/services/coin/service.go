package coin

import (
	"context"
	"fmt"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
	"github.com/nicolas-martin/memecoin-trading/pkg/dexscreens"
)

type Service struct {
	cache      redis.Cache
	dexScreens *dexscreens.Client
}

func NewService(cache redis.Cache, dexScreens *dexscreens.Client) *Service {
	return &Service{
		cache:      cache,
		dexScreens: dexScreens,
	}
}

func (s *Service) GetTopCoins(ctx context.Context, limit int) ([]models.Coin, error) {
	// Try to get from cache first
	coins, err := s.cache.GetTopCoins(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if coins != nil {
		return coins, nil
	}

	// If not in cache, fetch from DexScreens
	coins, err = s.dexScreens.GetTopCoins(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("dexscreens error: %w", err)
	}

	// Update cache in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cache.SetTopCoins(ctx, coins); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to cache top coins: %v\n", err)
		}
	}()

	return coins, nil
}

func (s *Service) GetCoinByID(ctx context.Context, id string) (*models.Coin, error) {
	// Try to get from cache first
	coin, err := s.cache.GetCoinByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if coin != nil {
		return coin, nil
	}

	// If not in cache, fetch from DexScreens
	coin, err = s.dexScreens.GetCoinByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("dexscreens error: %w", err)
	}

	// Update cache in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cache.SetCoin(ctx, coin); err != nil {
			fmt.Printf("Failed to cache coin: %v\n", err)
		}
	}()

	return coin, nil
}
