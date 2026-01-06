package redis

import (
	"context"
	"time"
	"wodge/internal/services"

	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client *redis.Client
}

func NewRedisDriver(addr string, password string, db int) (*RedisDriver, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisDriver{client: rdb}, nil
}

// Ensure RedisDriver implements services.CacheService
var _ services.CacheService = (*RedisDriver)(nil)

func (r *RedisDriver) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisDriver) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisDriver) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
