package abi

import (
	"context"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	redis_cache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type AbiManager struct {
	cache *Cache
	api   *eos.API
	ctx   context.Context
}

func NewAbiManager(rdb *redis.Client, api *eos.API, id string) *AbiManager {
	// Init abi cache
	cache := NewCache("ship.cache."+id+".abi", &redis_cache.Options{
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

func (mgr *AbiManager) GetAbi(account eos.AccountName) (*eos.ABI, error) {
	key := string(account)

	abi, err := mgr.cache.Get(key)
	if err != nil {
		resp, err := mgr.api.GetABI(mgr.ctx, account)
		if err != nil {
			return nil, fmt.Errorf("api: %s", err)
		}
		abi = &resp.ABI

		err = mgr.cache.Set(key, abi, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("cache: %s", err)
		}
	}
	return abi, nil
}
