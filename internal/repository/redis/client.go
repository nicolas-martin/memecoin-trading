package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nicolas-martin/memecoin-trading/internal/config"
)

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		// Pool configuration
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	return client, nil
}
