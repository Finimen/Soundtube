package repositories

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(keys...).Err()
}

func (r *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(keys...).Result()
}
