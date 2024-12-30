package memecoin

import (
	"context"
	"fmt"
	"log"
	"meme-trader/internal/repository/postgres"
	"time"
)

type Service struct {
	db        *postgres.Database
	providers []Provider
	logger    *log.Logger
}

func NewService(db *postgres.Database, logger *log.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
		providers: []Provider{
			NewDexScreenerProvider(),
			NewCoinGeckoProvider(),
			NewJupiterProvider(),
		},
	}
}

// GetTopMemeCoins returns the top meme coins by market cap
func (s *Service) GetTopMemeCoins(ctx context.Context, limit int) ([]postgres.MemeCoin, error) {
	return s.db.GetTopMemeCoins(limit)
}

// GetMemeCoinDetail returns detailed information about a specific meme coin
func (s *Service) GetMemeCoinDetail(ctx context.Context, coinID string) (*postgres.MemeCoin, []postgres.PriceHistory, error) {
	coin, err := s.db.GetMemeCoinByID(coinID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get meme coin: %w", err)
	}

	history, err := s.db.GetPriceHistory(coinID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get price history: %w", err)
	}

	return coin, history, nil
}

// FetchAndUpdateMemeCoins fetches meme coins from all providers and updates the database
func (s *Service) FetchAndUpdateMemeCoins(ctx context.Context) error {
	var memeCoins []postgres.MemeCoin

	for _, provider := range s.providers {
		coins, err := provider.FetchMemeCoins(ctx)
		if err != nil {
			s.logger.Printf("Error fetching meme coins from %s: %v", provider.Name(), err)
			continue
		}

		if len(coins) > 0 {
			memeCoins = append(memeCoins, coins...)
			break
		}
	}

	if len(memeCoins) == 0 {
		return fmt.Errorf("no meme coins found from any provider")
	}

	// Update database
	for _, coin := range memeCoins {
		if err := s.db.UpdateMemeCoin(&coin); err != nil {
			return fmt.Errorf("failed to update meme coin %s: %w", coin.Symbol, err)
		}

		// Add price history
		history := &postgres.PriceHistory{
			CoinID:    coin.ID,
			Price:     coin.Price,
			Volume:    coin.Volume24h,
			Timestamp: time.Now().Unix(),
		}

		if err := s.db.AddPriceHistory(history); err != nil {
			return fmt.Errorf("failed to add price history for %s: %w", coin.Symbol, err)
		}
	}

	return nil
}
