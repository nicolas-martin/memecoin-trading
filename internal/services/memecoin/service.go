package memecoin

import (
	"encoding/json"
	"fmt"
	"meme-trader/internal/repository/postgres"
	"net/http"
	"time"
)

type Service struct {
	db *postgres.Database
}

type DexScreensResponse struct {
	SchemaVersion string `json:"schemaVersion"`
	Pairs         []struct {
		ChainId   string `json:"chainId"`
		DexId     string `json:"dexId"`
		BaseToken struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"baseToken"`
		PriceUsd float64 `json:"priceUsd,string"`
		Volume   struct {
			H24 float64 `json:"h24"`
		} `json:"volume"`
		PriceChange struct {
			H24 float64 `json:"h24"`
		} `json:"priceChange"`
		Liquidity struct {
			Usd float64 `json:"usd"`
		} `json:"liquidity"`
		FDV       float64 `json:"fdv"`
		MarketCap float64 `json:"marketCap"`
	} `json:"pairs"`
}

func NewService(db *postgres.Database) *Service {
	return &Service{db: db}
}

func (s *Service) FetchAndUpdateMemeCoins() error {
	// Fetch data from dexscreens.io API
	resp, err := http.Get("https://api.dexscreener.com/latest/dex/search?q=solana")
	if err != nil {
		return fmt.Errorf("failed to fetch from dexscreener: %w", err)
	}
	defer resp.Body.Close()

	var dexResp DexScreensResponse
	if err := json.NewDecoder(resp.Body).Decode(&dexResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Update database with new information
	for _, pair := range dexResp.Pairs {
		if pair.ChainId != "solana" {
			continue // Skip non-Solana pairs
		}

		coin := &postgres.MemeCoin{
			ID:                       pair.BaseToken.Address,
			Symbol:                   pair.BaseToken.Symbol,
			Name:                     pair.BaseToken.Name,
			Price:                    pair.PriceUsd,
			MarketCap:                pair.MarketCap,
			Volume24h:                pair.Volume.H24,
			PriceChange24h:           pair.PriceChange.H24,
			PriceChangePercentage24h: pair.PriceChange.H24,
			ContractAddress:          pair.BaseToken.Address,
		}

		if err := s.db.UpdateMemeCoin(coin); err != nil {
			return fmt.Errorf("failed to update coin %s: %w", coin.Symbol, err)
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
