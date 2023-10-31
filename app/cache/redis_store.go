package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

type RedisStore struct {
	c   *cache.Cache
	ctx context.Context
}

func NewRedisStore(options *cache.Options) *RedisStore {
	return &RedisStore{
		c:   cache.New(options),
		ctx: context.Background(),
	}
}

func (s *RedisStore) Get(key string, value interface{}) error {
	return s.c.Get(s.ctx, key, value)
}

func (s *RedisStore) Has(key string) bool {
	return s.c.Exists(s.ctx, key)
}

func (s *RedisStore) Set(key string, value any, ttl time.Duration) error {
	return s.c.Set(&cache.Item{
		Ctx:   s.ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}
