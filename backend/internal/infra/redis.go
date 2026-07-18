package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	url := os.Getenv("URL_REDIS")
	if url == "" {
		url = "redis://localhost:6379/0"
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(fmt.Errorf("invalid redis URL %q: %w", url, err))
	}

	return &RedisClient{
		client: redis.NewClient(opts),
	}
}

func (r *RedisClient) Client() *redis.Client {
	return r.client
}

func (r *RedisClient) Ping(ctx context.Context) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return r.client.Ping(ctx).Err()
}
