package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

type RedisStore struct {
	c *cache.Cache
}

func NewRedisStore(options *cache.Options) *RedisStore {
	return &RedisStore{
		c: cache.New(options),
	}
}

func (s *RedisStore) Get(ctx context.Context, key string, value interface{}) error {
	return s.c.Get(ctx, key, value)
}

func (s *RedisStore) Has(ctx context.Context, key string) bool {
	return s.c.Exists(ctx, key)
}

func (s *RedisStore) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return s.c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}
