package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/karlseguin/typed"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	c *cache.Cache
}

type options struct {
	Stats bool
	Size  int
	TTL   time.Duration
}

func NewRedisStore(options *cache.Options) *RedisStore {
	return &RedisStore{
		c: cache.New(options),
	}
}

func getOptions(opts typed.Typed) options {
	return options{
		Stats: opts.Bool("stats"),
		Size:  opts.IntOr("size", 1000),
		TTL:   time.Duration(opts.IntOr("ttl", 10)) * time.Minute,
	}
}

func NewRedisFactory(client *redis.Client) Factory {
	return func(opts typed.Typed) (Store, error) {
		o := getOptions(opts)

		return NewRedisStore(&cache.Options{
			Redis:        client,
			StatsEnabled: o.Stats,
			LocalCache:   cache.NewTinyLFU(o.Size, o.TTL),
		}), nil
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
