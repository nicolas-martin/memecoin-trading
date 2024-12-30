package memecoin

import (
	"fmt"
	"meme-trader/internal/repository/postgres"
	"time"
)

type Service struct {
	db        *postgres.Database
	providers []Provider
}

func NewService(db *postgres.Database) *Service {
	return &Service{
		db: db,
		providers: []Provider{
			NewDexScreenerProvider(),
			NewCoinGeckoProvider(),
			NewJupiterProvider(),
		},
	}
}

func (s *Service) FetchAndUpdateMemeCoins() error {
	var lastErr error
	var coins []postgres.MemeCoin
	var successfulProvider Provider

	// Try each provider until we get data
	for _, provider := range s.providers {
		fetchedCoins, err := provider.FetchMemeCoins()
		if err == nil && len(fetchedCoins) > 0 {
			coins = fetchedCoins
			successfulProvider = provider
			break
		}
		lastErr = err
	}

	if lastErr != nil && len(coins) == 0 {
		return fmt.Errorf("all providers failed, last error: %w", lastErr)
	}

	if len(coins) == 0 {
		return fmt.Errorf("no coins found from any provider")
	}

	// Update database with new information
	for i := range coins {
		// Add provider information
		coins[i].DataProvider = successfulProvider.Name()
		coins[i].LastUpdated = time.Now()

		if err := s.db.UpdateMemeCoin(&coins[i]); err != nil {
			return fmt.Errorf("failed to update coin %s: %w", coins[i].Symbol, err)
		}

		// Add price history
		history := &postgres.PriceHistory{
			CoinID:    coins[i].ID,
			Price:     coins[i].Price,
			Volume:    coins[i].Volume24h,
			Timestamp: time.Now().Unix(),
		}

		if err := s.db.AddPriceHistory(history); err != nil {
			return fmt.Errorf("failed to add price history for %s: %w", coins[i].Symbol, err)
		}
	}

	return nil
}

func (s *Service) GetTopMemeCoins(limit int) ([]postgres.MemeCoin, error) {
	return s.db.GetTopMemeCoins(limit)
}

func (s *Service) GetMemeCoinDetail(id string) (*postgres.MemeCoin, []postgres.PriceHistory, error) {
	coin, err := s.db.GetMemeCoinByID(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get coin: %w", err)
	}

	history, err := s.db.GetPriceHistory(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get price history: %w", err)
	}

	return coin, history, nil
}
