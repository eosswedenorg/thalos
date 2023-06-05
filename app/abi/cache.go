package abi

import (
	"context"
	"time"

	eos "github.com/eoscanada/eos-go"
	redis_cache "github.com/go-redis/cache/v9"
)

type Cache struct {
	c      *redis_cache.Cache
	ctx    context.Context
	prefix string
}

func NewCache(prefix string, options *redis_cache.Options) *Cache {
	return &Cache{
		c:      redis_cache.New(options),
		ctx:    context.Background(),
		prefix: prefix,
	}
}

func (cache *Cache) Get(account string) (*eos.ABI, error) {
	var v eos.ABI
	err := cache.c.Get(cache.ctx, cache.key(account), &v)
	return &v, err
}

func (cache *Cache) Set(account string, abi *eos.ABI, ttl time.Duration) error {
	return cache.c.Set(&redis_cache.Item{
		Ctx:   cache.ctx,
		Key:   cache.key(account),
		Value: *abi,
		TTL:   ttl,
	})
}

func (cache *Cache) key(account string) string {
	return cache.prefix + "." + account
}
