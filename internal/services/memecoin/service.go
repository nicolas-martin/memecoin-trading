package memecoin

import (
	"context"
	"fmt"
	"log"
	"meme-trader/internal/repository/postgres"
	"os"
	"time"
)

type Service struct {
	db        *postgres.Database
	providers []Provider
	logger    *log.Logger
}

func NewService(db *postgres.Database, logger *log.Logger) *Service {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
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
	s.logger.Println("Starting to fetch meme coins from all providers")
	var memeCoins []postgres.MemeCoin
	coinMap := make(map[string]*postgres.MemeCoin) // Map to track coins by address

	for _, provider := range s.providers {
		s.logger.Printf("Fetching from provider: %s", provider.Name())
		coins, err := provider.FetchMemeCoins(ctx)
		if err != nil {
			s.logger.Printf("Error fetching meme coins from %s: %v", provider.Name(), err)
			continue
		}

		s.logger.Printf("Got %d coins from %s", len(coins), provider.Name())

		// Merge data from this provider
		for _, coin := range coins {
			if existing, ok := coinMap[coin.ContractAddress]; ok {
				// Update existing coin with non-empty values
				if coin.LogoURL != "" {
					s.logger.Printf("Updating logo URL for %s from %s to %s", coin.Symbol, existing.LogoURL, coin.LogoURL)
					existing.LogoURL = coin.LogoURL
				}
				if coin.Name != "" {
					existing.Name = coin.Name
				}
				if coin.Symbol != "" {
					existing.Symbol = coin.Symbol
				}
				if coin.Price > 0 {
					existing.Price = coin.Price
				}
				if coin.MarketCap > 0 {
					existing.MarketCap = coin.MarketCap
				}
				if coin.Volume24h > 0 {
					existing.Volume24h = coin.Volume24h
				}
				existing.LastUpdated = time.Now()
			} else {
				// Add new coin to map
				s.logger.Printf("Adding new coin %s with logo URL: %s", coin.Symbol, coin.LogoURL)
				coinCopy := coin
				coinMap[coin.ContractAddress] = &coinCopy
			}
		}
	}

	// Convert map back to slice
	for _, coin := range coinMap {
		memeCoins = append(memeCoins, *coin)
	}

	if len(memeCoins) == 0 {
		return fmt.Errorf("no meme coins found from any provider")
	}

	s.logger.Printf("Updating database with %d meme coins", len(memeCoins))

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

	s.logger.Println("Successfully updated all meme coins")
	return nil
}
