package abi_cache

import (
	"context"
	"time"

	eos "github.com/eoscanada/eos-go"
	redis_cache "github.com/go-redis/cache/v8"
)

type Cache struct {
	c      *redis_cache.Cache
	ctx    context.Context
	prefix string
}

func New(prefix string, options *redis_cache.Options) *Cache {
	return &Cache{
		c:      redis_cache.New(options),
		ctx:    context.Background(),
		prefix: prefix,
	}
}

func (this *Cache) Get(account string) (*eos.ABI, error) {
	var v eos.ABI
	err := this.c.Get(this.ctx, this.key(account), &v)
	return &v, err
}

func (this *Cache) Set(account string, abi *eos.ABI, ttl time.Duration) error {
	return this.c.Set(&redis_cache.Item{
		Ctx:   this.ctx,
		Key:   this.key(account),
		Value: *abi,
		TTL:   ttl,
	})
}

func (this *Cache) key(account string) string {
	return this.prefix + "." + account
}
