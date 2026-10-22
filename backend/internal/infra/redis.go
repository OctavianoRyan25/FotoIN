package infra

import (
	"fmt"

	"fotoin/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	return &RedisClient{
		client: redis.NewClient(opts),
	}, nil
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}
