package cache

import (
	"time"
)

type Cache struct {
	store  Store
	prefix string
}

// Create a new cache
func NewCache(prefix string, store Store) *Cache {
	return &Cache{
		store:  store,
		prefix: prefix,
	}
}

func (cache *Cache) Get(key string, value any) error {
	return cache.store.Get(cache.key(key), value)
}

func (cache *Cache) Set(key string, value any, ttl time.Duration) error {
	return cache.store.Set(cache.key(key), value, ttl)
}

func (cache *Cache) key(key string) string {
	return cache.prefix + "::" + key
}
