package memecoin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"meme-trader/internal/repository/postgres"
	"net/http"
	"strings"
	"time"
)

// Provider interface for fetching meme coins
type Provider interface {
	Name() string
	FetchMemeCoins(ctx context.Context) ([]postgres.MemeCoin, error)
}

// DexScreenerProvider implements the Provider interface for DexScreener
type DexScreenerProvider struct {
	client *http.Client
}

func NewDexScreenerProvider() *DexScreenerProvider {
	return &DexScreenerProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *DexScreenerProvider) Name() string {
	return "DexScreener"
}

func (p *DexScreenerProvider) FetchMemeCoins(ctx context.Context) ([]postgres.MemeCoin, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.dexscreener.com/latest/dex/search?q=solana", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response struct {
		Pairs []struct {
			ChainID     string `json:"chainId"`
			DexID       string `json:"dexId"`
			URL         string `json:"url"`
			PairAddress string `json:"pairAddress"`
			BaseToken   struct {
				Address     string `json:"address"`
				Name        string `json:"name"`
				Symbol      string `json:"symbol"`
				LogoURL     string `json:"logoURI"`
				PriceUSD    string `json:"priceUSD"`
				PriceChange struct {
					H24 float64 `json:"h24"`
				} `json:"priceChange"`
			} `json:"baseToken"`
			Volume struct {
				H24 float64 `json:"h24"`
			} `json:"volume"`
		} `json:"pairs"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var coins []postgres.MemeCoin
	seen := make(map[string]bool)

	for _, pair := range response.Pairs {
		if pair.ChainID != "solana" {
			continue
		}

		if seen[pair.BaseToken.Address] {
			continue
		}

		price := 0.0
		if pair.BaseToken.PriceUSD != "" {
			fmt.Sscanf(pair.BaseToken.PriceUSD, "%f", &price)
		}

		coin := postgres.MemeCoin{
			ID:                       pair.BaseToken.Address,
			Symbol:                   pair.BaseToken.Symbol,
			Name:                     pair.BaseToken.Name,
			Price:                    price,
			Volume24h:                pair.Volume.H24,
			PriceChangePercentage24h: pair.BaseToken.PriceChange.H24,
			ContractAddress:          pair.BaseToken.Address,
			DataProvider:             "DexScreener",
			LastUpdated:              time.Now(),
			LogoURL:                  pair.BaseToken.LogoURL,
		}

		coins = append(coins, coin)
		seen[pair.BaseToken.Address] = true
	}

	return coins, nil
}

// CoinGeckoProvider implements the Provider interface for CoinGecko
type CoinGeckoProvider struct {
	client *http.Client
}

func NewCoinGeckoProvider() *CoinGeckoProvider {
	return &CoinGeckoProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *CoinGeckoProvider) Name() string {
	return "CoinGecko"
}

func (p *CoinGeckoProvider) FetchMemeCoins(ctx context.Context) ([]postgres.MemeCoin, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&category=meme-token&order=market_cap_desc&per_page=100&page=1&sparkline=false", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	var response []struct {
		ID                       string  `json:"id"`
		Symbol                   string  `json:"symbol"`
		Name                     string  `json:"name"`
		Image                    string  `json:"image"`
		CurrentPrice             float64 `json:"current_price"`
		MarketCap                float64 `json:"market_cap"`
		TotalVolume              float64 `json:"total_volume"`
		PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
		PriceChange24h           float64 `json:"price_change_24h"`
		ContractAddress          string  `json:"contract_address"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var coins []postgres.MemeCoin
	for _, item := range response {
		coin := postgres.MemeCoin{
			ID:                       item.ID,
			Symbol:                   strings.ToUpper(item.Symbol),
			Name:                     item.Name,
			Price:                    item.CurrentPrice,
			MarketCap:                item.MarketCap,
			Volume24h:                item.TotalVolume,
			PriceChange24h:           item.PriceChange24h,
			PriceChangePercentage24h: item.PriceChangePercentage24h,
			ContractAddress:          item.ContractAddress,
			DataProvider:             "CoinGecko",
			LastUpdated:              time.Now(),
			LogoURL:                  item.Image,
		}
		coins = append(coins, coin)
	}

	return coins, nil
}

// JupiterProvider implements the Provider interface for Jupiter
type JupiterProvider struct {
	client *http.Client
}

func NewJupiterProvider() *JupiterProvider {
	return &JupiterProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *JupiterProvider) Name() string {
	return "Jupiter"
}

func (p *JupiterProvider) FetchMemeCoins(ctx context.Context) ([]postgres.MemeCoin, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://token.jup.ag/all", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Tokens []struct {
			Address   string   `json:"address"`
			Symbol    string   `json:"symbol"`
			Name      string   `json:"name"`
			LogoURI   string   `json:"logoURI"`
			Price     float64  `json:"price"`
			Volume24h float64  `json:"volume24h"`
			MarketCap float64  `json:"marketCap"`
			Tags      []string `json:"tags"`
		} `json:"tokens"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var coins []postgres.MemeCoin
	for _, token := range response.Tokens {
		// Filter for meme tokens
		isMeme := false
		for _, tag := range token.Tags {
			if strings.Contains(strings.ToLower(tag), "meme") {
				isMeme = true
				break
			}
		}
		if !isMeme {
			continue
		}

		coin := postgres.MemeCoin{
			ID:              token.Address,
			Symbol:          token.Symbol,
			Name:            token.Name,
			Price:           token.Price,
			MarketCap:       token.MarketCap,
			Volume24h:       token.Volume24h,
			ContractAddress: token.Address,
			DataProvider:    "Jupiter",
			LastUpdated:     time.Now(),
			LogoURL:         token.LogoURI,
		}
		coins = append(coins, coin)
	}

	return coins, nil
}
