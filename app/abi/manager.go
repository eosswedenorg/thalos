package abi

import (
	"context"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	redis_cache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// AbiManager handles an ABI cache that fetches the ABI from an API on cache miss.
type AbiManager struct {
	cache *Cache
	api   *eos.API
	ctx   context.Context
}

// Create a new ABI Manager
func NewAbiManager(rdb *redis.Client, api *eos.API, id string) *AbiManager {
	// Init abi cache
	cache := NewCache("thalos::cache::"+id+"::abi", &redis_cache.Options{
		Redis: rdb,
		// Cache 10k keys for 10 minutes.
		LocalCache: redis_cache.NewTinyLFU(10000, 10*time.Minute),
	})

	return &AbiManager{
		cache: cache,
		api:   api,
		ctx:   context.Background(),
	}
}

// Set or update an ABI in the cache.
func (mgr *AbiManager) SetAbi(account eos.AccountName, abi *eos.ABI) error {
	return mgr.cache.Set(string(account), abi, time.Hour)
}

// Get an ABI from the cache, on cache miss it is fetched from the
// API, gets cached and then returned to the user
func (mgr *AbiManager) GetAbi(account eos.AccountName) (*eos.ABI, error) {
	key := string(account)

	abi, err := mgr.cache.Get(key)
	if err != nil {
		resp, err := mgr.api.GetABI(mgr.ctx, account)
		if err != nil {
			return nil, fmt.Errorf("api: %s", err)
		}
		abi = &resp.ABI

		err = mgr.SetAbi(account, abi)
		if err != nil {
			return nil, fmt.Errorf("cache: %s", err)
		}
	}
	return abi, nil
}
