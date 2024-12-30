package memecoin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// First fetch the token boosts to get logos
	boostReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.dexscreener.com/token-boosts/latest/v1", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create boost request: %w", err)
	}

	boostResp, err := p.client.Do(boostReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch boost data: %w", err)
	}
	defer boostResp.Body.Close()

	var boostResponse struct {
		TokenAddress string `json:"tokenAddress"`
		Icon         string `json:"icon"`
	}

	if err := json.NewDecoder(boostResp.Body).Decode(&boostResponse); err != nil {
		return nil, fmt.Errorf("failed to parse boost response: %w", err)
	}

	// Create a map of token addresses to logos
	logoMap := make(map[string]string)
	if boostResponse.TokenAddress != "" && boostResponse.Icon != "" {
		logoMap[strings.ToLower(boostResponse.TokenAddress)] = boostResponse.Icon
	}

	// Now fetch the main token data
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

		// Try to get logo from boost data first, fallback to token data
		logoURL := pair.BaseToken.LogoURL
		if boostLogo, ok := logoMap[strings.ToLower(pair.BaseToken.Address)]; ok && boostLogo != "" {
			logoURL = boostLogo
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
			LogoURL:                  logoURL,
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

	log.Printf("CoinGecko: Found %d meme tokens", len(response))

	var coins []postgres.MemeCoin
	for _, item := range response {
		log.Printf("CoinGecko: Processing token %s (%s) with logo: %s", item.Name, item.Symbol, item.Image)

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

	log.Printf("CoinGecko: Processed %d meme tokens", len(coins))
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

	var tokens []struct {
		Address   string   `json:"address"`
		Symbol    string   `json:"symbol"`
		Name      string   `json:"name"`
		LogoURI   string   `json:"logoURI"`
		Price     float64  `json:"price"`
		Volume24h float64  `json:"volume24h"`
		MarketCap float64  `json:"marketCap"`
		Tags      []string `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("Jupiter: Found %d total tokens", len(tokens))

	var coins []postgres.MemeCoin
	memeCount := 0
	for _, token := range tokens {
		// Filter for meme tokens using more lenient criteria
		isMeme := false
		tokenNameLower := strings.ToLower(token.Name)
		tokenSymbolLower := strings.ToLower(token.Symbol)

		// Check tags first
		for _, tag := range token.Tags {
			if strings.Contains(strings.ToLower(tag), "meme") {
				isMeme = true
				break
			}
		}

		// If not found in tags, check name and symbol for common meme indicators
		if !isMeme {
			memeIndicators := []string{
				"doge", "shib", "pepe", "wojak", "chad", "inu", "cat",
				"moon", "elon", "safe", "baby", "rocket", "meme", "bonk",
				"floki", "cheems", "frog", "ape", "monkey", "dog", "wow",
				"based", "wagmi", "gm", "chad", "whale", "bear", "bull",
				"wojak", "nft", "moon", "lambo", "tendies", "diamond",
				"hands", "hodl", "fomo", "yolo", "wen", "ser", "ngmi",
			}

			for _, indicator := range memeIndicators {
				if strings.Contains(tokenNameLower, indicator) || strings.Contains(tokenSymbolLower, indicator) {
					isMeme = true
					break
				}
			}
		}

		if !isMeme {
			continue
		}
		memeCount++

		// Log token details for debugging
		log.Printf("Jupiter: Processing meme token: %s (%s)", token.Name, token.Address)
		log.Printf("Jupiter: Original LogoURI: %s", token.LogoURI)

		// Ensure we have a valid logo URL
		logoURL := token.LogoURI
		if !strings.HasPrefix(logoURL, "http") && !strings.HasPrefix(logoURL, "https") {
			// Try to construct a valid URL if it's a relative path
			if strings.HasPrefix(logoURL, "/") {
				logoURL = "https://token.jup.ag" + logoURL
				log.Printf("Jupiter: Converted relative path to: %s", logoURL)
			} else {
				// If no logo URL is provided, try to get it from another source
				logoURL = fmt.Sprintf("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/%s/logo.png", token.Address)
				log.Printf("Jupiter: Using fallback logo URL: %s", logoURL)
			}
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
			LogoURL:         logoURL,
		}
		coins = append(coins, coin)
	}

	log.Printf("Jupiter: Found %d meme tokens out of %d total tokens", memeCount, len(tokens))
	return coins, nil
}
