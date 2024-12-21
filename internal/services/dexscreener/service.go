package dexscreener

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicolas-martin/memecoin-trading/internal/repository/redis"
)

type Service struct {
	httpClient *http.Client
	cache      redis.Cache
	baseURL    string
}

type PriceData struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

type DexScreenerResponse struct {
	Data struct {
		Prices []struct {
			Timestamp int64   `json:"timestamp"`
			PriceUsd  float64 `json:"priceUsd,string"`
		} `json:"prices"`
	} `json:"data"`
}

func NewService(cache redis.Cache) *Service {
	return &Service{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      cache,
		baseURL:    "https://api.dexscreener.com/latest",
	}
}

func (s *Service) GetHistoricalPrices(ctx context.Context, pairAddress string, timeframe string) ([]PriceData, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("prices:%s:%s", pairAddress, timeframe)
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var prices []PriceData
		if err := json.Unmarshal([]byte(cachedData), &prices); err == nil {
			return prices, nil
		}
	}

	// If not in cache or error, fetch from DexScreener
	prices, err := s.fetchFromDexScreener(ctx, pairAddress, timeframe)
	if err != nil {
		return nil, err
	}

	// Cache the results
	cacheDuration := getCacheDuration(timeframe)
	data, _ := json.Marshal(prices)
	s.cache.Set(ctx, cacheKey, string(data), cacheDuration)

	return prices, nil
}

func (s *Service) fetchFromDexScreener(ctx context.Context, pairAddress string, timeframe string) ([]PriceData, error) {
	from := getFromTimestamp(timeframe)
	url := fmt.Sprintf("%s/dex/pairs/%s/prices?from=%d", s.baseURL, pairAddress, from)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response DexScreenerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	prices := make([]PriceData, len(response.Data.Prices))
	for i, p := range response.Data.Prices {
		prices[i] = PriceData{
			Timestamp: p.Timestamp,
			Price:     p.PriceUsd,
		}
	}

	return prices, nil
}

func getFromTimestamp(timeframe string) int64 {
	now := time.Now().Unix() * 1000
	switch timeframe {
	case "1H":
		return now - 60*60*1000
	case "24H":
		return now - 24*60*60*1000
	case "1W":
		return now - 7*24*60*60*1000
	case "1M":
		return now - 30*24*60*60*1000
	case "1Y":
		return now - 365*24*60*60*1000
	default:
		return now - 24*60*60*1000
	}
}

func getCacheDuration(timeframe string) time.Duration {
	switch timeframe {
	case "1H":
		return 1 * time.Minute
	case "24H":
		return 5 * time.Minute
	case "1W":
		return 15 * time.Minute
	case "1M":
		return 1 * time.Hour
	case "1Y":
		return 24 * time.Hour
	default:
		return 5 * time.Minute
	}
}
