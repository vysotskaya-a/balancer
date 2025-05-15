package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
}

func New(addr string, db int, logger *zap.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client, logger: logger}, nil
}

func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		c.logger.Error("Redis Get failed", zap.Error(err))
		return "", err
	}
	return val, nil
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), key, jsonData, ttl).Err()
}
