package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
)

const (
	// Key prefixes
	coinPrefix           = "coin:"
	userPrefix           = "user:"
	statsPrefix          = "stats:"
	topCoinsKey          = "top_coins"
	defaultTTL           = 15 * time.Minute
	longTTL              = 1 * time.Hour
	statsShortTTL        = 5 * time.Minute
	leaderboardKeyPrefix = "leaderboard:"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	GetTopCoins(ctx context.Context, limit int) ([]models.Coin, error)
	SetTopCoins(ctx context.Context, coins []models.Coin) error
	GetCoinByID(ctx context.Context, id uuid.UUID) (*models.Coin, error)
	SetCoin(ctx context.Context, coin *models.Coin) error
	InvalidateCoinCache(ctx context.Context, id uuid.UUID) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	SetUser(ctx context.Context, user *models.User) error
	InvalidateUserCache(ctx context.Context, id string) error
	GetUserStats(ctx context.Context, userID string) (map[string]float64, error)
	SetUserStats(ctx context.Context, userID string, stats map[string]float64) error
	InvalidateUserStats(ctx context.Context, userID string) error
	GetLeaderboard(ctx context.Context, timeframe string) ([]models.LeaderboardEntry, error)
	SetLeaderboard(ctx context.Context, timeframe string, entries []models.LeaderboardEntry) error
	InvalidateLeaderboard(ctx context.Context, timeframe string) error
	Clear(ctx context.Context) error
	Ping(ctx context.Context) error
}

// RedisCache implements the Cache interface
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Basic operations
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Coin methods
func (c *RedisCache) GetTopCoins(ctx context.Context, limit int) ([]models.Coin, error) {
	key := fmt.Sprintf("%s:%d", topCoinsKey, limit)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var coins []models.Coin
	if err := json.Unmarshal(data, &coins); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return coins, nil
}

func (c *RedisCache) SetTopCoins(ctx context.Context, coins []models.Coin) error {
	data, err := json.Marshal(coins)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	key := fmt.Sprintf("%s:%d", topCoinsKey, len(coins))
	if err := c.client.Set(ctx, key, data, defaultTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (c *RedisCache) GetCoinByID(ctx context.Context, id uuid.UUID) (*models.Coin, error) {
	key := coinPrefix + id.String()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var coin models.Coin
	if err := json.Unmarshal(data, &coin); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &coin, nil
}

func (c *RedisCache) SetCoin(ctx context.Context, coin *models.Coin) error {
	data, err := json.Marshal(coin)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	key := coinPrefix + coin.ID.String()
	if err := c.client.Set(ctx, key, data, longTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (c *RedisCache) InvalidateCoinCache(ctx context.Context, id uuid.UUID) error {
	key := coinPrefix + id.String()
	return c.client.Del(ctx, key).Err()
}

// User methods
func (c *RedisCache) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	key := userPrefix + id
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &user, nil
}

func (c *RedisCache) SetUser(ctx context.Context, user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	key := userPrefix + user.ID.String()
	if err := c.client.Set(ctx, key, data, longTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (c *RedisCache) InvalidateUserCache(ctx context.Context, id string) error {
	key := userPrefix + id
	return c.client.Del(ctx, key).Err()
}

// Transaction stats methods
func (c *RedisCache) GetUserStats(ctx context.Context, userID string) (map[string]float64, error) {
	key := statsPrefix + userID
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var stats map[string]float64
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return stats, nil
}

func (c *RedisCache) SetUserStats(ctx context.Context, userID string, stats map[string]float64) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	key := statsPrefix + userID
	if err := c.client.Set(ctx, key, data, statsShortTTL).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

func (c *RedisCache) InvalidateUserStats(ctx context.Context, userID string) error {
	key := statsPrefix + userID
	return c.client.Del(ctx, key).Err()
}

// Helper methods
func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Leaderboard methods
func (c *RedisCache) GetLeaderboard(ctx context.Context, timeframe string) ([]models.LeaderboardEntry, error) {
	key := leaderboardKeyPrefix + timeframe
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var entries []models.LeaderboardEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *RedisCache) SetLeaderboard(ctx context.Context, timeframe string, entries []models.LeaderboardEntry) error {
	key := leaderboardKeyPrefix + timeframe
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	// Cache leaderboard for 5 minutes
	return c.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (c *RedisCache) InvalidateLeaderboard(ctx context.Context, timeframe string) error {
	key := leaderboardKeyPrefix + timeframe
	return c.client.Del(ctx, key).Err()
}
