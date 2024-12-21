package coin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
)

const (
	defaultTTL  = 15 * time.Minute
	trendingTTL = 5 * time.Minute
	longTTL     = 1 * time.Hour
)

type Service struct {
	cache      redis.Cache
	httpClient *http.Client
	baseURL    string
}

type DexScreenerResponse struct {
	SchemaVersion string `json:"schemaVersion"`
	Pairs         []struct {
		ChainId     string `json:"chainId"`
		DexId       string `json:"dexId"`
		URL         string `json:"url"`
		PairAddress string `json:"pairAddress"`
		BaseToken   struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"baseToken"`
		QuoteToken struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"quoteToken"`
		PriceNative string `json:"priceNative"`
		PriceUsd    string `json:"priceUsd"`
		Liquidity   struct {
			USD   float64 `json:"usd"`
			Base  float64 `json:"base"`
			Quote float64 `json:"quote"`
		} `json:"liquidity"`
		Volume struct {
			H24 float64 `json:"h24"`
			H6  float64 `json:"h6"`
			H1  float64 `json:"h1"`
			M5  float64 `json:"m5"`
		} `json:"volume"`
		PriceChange struct {
			H1  float64 `json:"h1"`
			H24 float64 `json:"h24"`
			D7  float64 `json:"d7"`
		} `json:"priceChange"`
		FDV       float64 `json:"fdv"`
		MarketCap float64 `json:"marketCap"`
		Info      struct {
			ImageURL    string `json:"imageUrl,omitempty"`
			Description string `json:"description,omitempty"`
			Websites    []struct {
				URL string `json:"url"`
			} `json:"websites,omitempty"`
			Socials []struct {
				Platform string `json:"platform"`
				Handle   string `json:"handle"`
			} `json:"socials,omitempty"`
		} `json:"info,omitempty"`
	} `json:"pairs"`
}

type PricePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

func (s *Service) GetPairData(ctx context.Context, chainId, pairAddress string) (*models.Coin, error) {
	cacheKey := fmt.Sprintf("pair:%s:%s", chainId, pairAddress)

	// Try cache first
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var coin models.Coin
		if err := json.Unmarshal([]byte(cached), &coin); err == nil {
			return &coin, nil
		}
	}

	// Fetch from DexScreener
	url := fmt.Sprintf("%s/pairs/%s/%s", s.baseURL, chainId, pairAddress)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dexResp DexScreenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dexResp); err != nil {
		return nil, err
	}

	if len(dexResp.Pairs) == 0 {
		return nil, fmt.Errorf("pair not found")
	}

	pair := dexResp.Pairs[0]
	coin := &models.Coin{
		Name:        pair.BaseToken.Name,
		Symbol:      pair.BaseToken.Symbol,
		PairAddress: pair.PairAddress,
		ChainID:     pair.ChainId,
		Price:       pair.PriceUsd,
		PriceChange: models.PriceChange{
			H1:  pair.PriceChange.H1,
			H24: pair.PriceChange.H24,
			D7:  pair.PriceChange.D7,
		},
		Volume: models.Volume{
			H24: pair.Volume.H24,
			H6:  pair.Volume.H6,
			H1:  pair.Volume.H1,
			M5:  pair.Volume.M5,
		},
		Liquidity: models.Liquidity{
			USD:   pair.Liquidity.USD,
			Base:  pair.Liquidity.Base,
			Quote: pair.Liquidity.Quote,
		},
		MarketCap:   pair.MarketCap,
		FDV:         pair.FDV,
		Logo:        pair.Info.ImageURL,
		Description: pair.Info.Description,
	}

	// Cache the result
	if data, err := json.Marshal(coin); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), defaultTTL)
	}

	return coin, nil
}

func (s *Service) GetTrendingCoins(ctx context.Context) ([]*models.Coin, error) {
	// Try cache first
	if cached, err := s.cache.Get(ctx, "trending_coins"); err == nil {
		var coins []*models.Coin
		if err := json.Unmarshal([]byte(cached), &coins); err == nil {
			return coins, nil
		}
	}

	// Fetch from DexScreener
	url := fmt.Sprintf("%s/trending", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coins []*models.Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, err
	}

	// Cache the result
	if data, err := json.Marshal(coins); err == nil {
		s.cache.Set(ctx, "trending_coins", string(data), trendingTTL)
	}

	return coins, nil
}

func (s *Service) GetTopCoins(ctx context.Context, limit int) ([]*models.Coin, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("top_coins:%d", limit)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
		var coins []*models.Coin
		if err := json.Unmarshal([]byte(cached), &coins); err == nil {
			return coins, nil
		}
	}

	// Fetch from DexScreener - using volume to determine top coins
	url := fmt.Sprintf("%s/search?q=volume>1000000 sort:volume.h24:desc&limit=%d", s.baseURL, limit)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dexResp DexScreenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&dexResp); err != nil {
		return nil, err
	}

	coins := make([]*models.Coin, 0, len(dexResp.Pairs))
	for _, pair := range dexResp.Pairs {
		coin := &models.Coin{
			Name:        pair.BaseToken.Name,
			Symbol:      pair.BaseToken.Symbol,
			PairAddress: pair.PairAddress,
			ChainID:     pair.ChainId,
			Price:       pair.PriceUsd,
			PriceChange: models.PriceChange{
				H1:  pair.PriceChange.H1,
				H24: pair.PriceChange.H24,
				D7:  pair.PriceChange.D7,
			},
			Volume: models.Volume{
				H24: pair.Volume.H24,
				H6:  pair.Volume.H6,
				H1:  pair.Volume.H1,
				M5:  pair.Volume.M5,
			},
			Liquidity: models.Liquidity{
				USD:   pair.Liquidity.USD,
				Base:  pair.Liquidity.Base,
				Quote: pair.Liquidity.Quote,
			},
			MarketCap:   pair.MarketCap,
			FDV:         pair.FDV,
			Logo:        pair.Info.ImageURL,
			Description: pair.Info.Description,
		}
		coins = append(coins, coin)
	}

	// Cache the result
	if data, err := json.Marshal(coins); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), defaultTTL)
	}

	return coins, nil
}

func (s *Service) GetHistoricalPrices(ctx context.Context, pairAddress string) ([]PricePoint, error) {
	cacheKey := fmt.Sprintf("historical_prices:%s", pairAddress)

	// Try to get from cache first
	if cachedData, err := s.cache.Get(ctx, cacheKey); err == nil {
		var prices []PricePoint
		if err := json.Unmarshal([]byte(cachedData), &prices); err == nil {
			return prices, nil
		}
	}

	// If not in cache, fetch from API
	url := fmt.Sprintf("%s/pairs/history/%s", s.baseURL, pairAddress)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	var prices []PricePoint
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Cache the results
	if jsonData, err := json.Marshal(prices); err == nil {
		if err := s.cache.Set(ctx, cacheKey, string(jsonData), longTTL); err != nil {
			// Log the error but don't fail the request
			fmt.Printf("failed to cache historical prices: %v\n", err)
		}
	}

	return prices, nil
}

func NewService(cache redis.Cache) *Service {
	return &Service{
		cache:      cache,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.dexscreener.com/latest/dex",
	}
}
