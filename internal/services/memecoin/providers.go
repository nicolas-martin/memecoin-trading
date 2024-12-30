package memecoin

import (
	"encoding/json"
	"fmt"
	"meme-trader/internal/repository/postgres"
	"net/http"
	"time"
)

// Provider defines the interface for meme coin data providers
type Provider interface {
	Name() string
	FetchMemeCoins() ([]postgres.MemeCoin, error)
}

// DexScreenerProvider implements the DexScreener API
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

func (p *DexScreenerProvider) FetchMemeCoins() ([]postgres.MemeCoin, error) {
	resp, err := p.client.Get("https://api.dexscreener.com/latest/dex/search?q=solana")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from dexscreener: %w", err)
	}
	defer resp.Body.Close()

	var dexResp DexScreensResponse
	if err := json.NewDecoder(resp.Body).Decode(&dexResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var coins []postgres.MemeCoin
	for _, pair := range dexResp.Pairs {
		if pair.ChainId != "solana" {
			continue
		}

		coin := postgres.MemeCoin{
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
		coins = append(coins, coin)
	}

	return coins, nil
}

// CoinGeckoProvider implements the CoinGecko API
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

type CoinGeckoResponse struct {
	Coins []struct {
		ID            string  `json:"id"`
		Symbol        string  `json:"symbol"`
		Name          string  `json:"name"`
		CurrentPrice  float64 `json:"current_price"`
		MarketCap     float64 `json:"market_cap"`
		Volume24h     float64 `json:"total_volume"`
		PriceChange24 float64 `json:"price_change_24h"`
		PriceChange   float64 `json:"price_change_percentage_24h"`
		Platforms     struct {
			Solana string `json:"solana"`
		} `json:"platforms"`
	} `json:"coins"`
}

func (p *CoinGeckoProvider) FetchMemeCoins() ([]postgres.MemeCoin, error) {
	resp, err := p.client.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&category=meme-token&order=market_cap_desc&per_page=100&page=1&sparkline=false")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from coingecko: %w", err)
	}
	defer resp.Body.Close()

	var geckoResp CoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geckoResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var coins []postgres.MemeCoin
	for _, coin := range geckoResp.Coins {
		if coin.Platforms.Solana == "" {
			continue
		}

		memeCoin := postgres.MemeCoin{
			ID:                       coin.ID,
			Symbol:                   coin.Symbol,
			Name:                     coin.Name,
			Price:                    coin.CurrentPrice,
			MarketCap:                coin.MarketCap,
			Volume24h:                coin.Volume24h,
			PriceChange24h:           coin.PriceChange24,
			PriceChangePercentage24h: coin.PriceChange,
			ContractAddress:          coin.Platforms.Solana,
		}
		coins = append(coins, memeCoin)
	}

	return coins, nil
}

// JupiterProvider implements the Jupiter API
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

type JupiterResponse struct {
	Data []struct {
		Address     string  `json:"address"`
		Symbol      string  `json:"symbol"`
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		MarketCap   float64 `json:"marketCap"`
		Volume24h   float64 `json:"volume24h"`
		PriceChange struct {
			Percentage24h float64 `json:"percentage24h"`
			Value24h      float64 `json:"value24h"`
		} `json:"priceChange"`
	} `json:"data"`
}

func (p *JupiterProvider) FetchMemeCoins() ([]postgres.MemeCoin, error) {
	resp, err := p.client.Get("https://price.jup.ag/v4/token-list")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from jupiter: %w", err)
	}
	defer resp.Body.Close()

	var jupResp JupiterResponse
	if err := json.NewDecoder(resp.Body).Decode(&jupResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var coins []postgres.MemeCoin
	for _, token := range jupResp.Data {
		coin := postgres.MemeCoin{
			ID:                       token.Address,
			Symbol:                   token.Symbol,
			Name:                     token.Name,
			Price:                    token.Price,
			MarketCap:                token.MarketCap,
			Volume24h:                token.Volume24h,
			PriceChange24h:           token.PriceChange.Value24h,
			PriceChangePercentage24h: token.PriceChange.Percentage24h,
			ContractAddress:          token.Address,
		}
		coins = append(coins, coin)
	}

	return coins, nil
}
