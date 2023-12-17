package abi

import (
	"context"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/eosswedenorg/thalos/app/cache"
)

// AbiManager handles an ABI cache that fetches the ABI from an API on cache miss.
type AbiManager struct {
	cache *cache.Cache
	api   *eos.API
	ctx   context.Context
}

// Create a new ABI Manager
func NewAbiManager(cache *cache.Cache, api *eos.API) *AbiManager {
	return &AbiManager{
		cache: cache,
		api:   api,
		ctx:   context.Background(),
	}
}

// Set or update an ABI in the cache.
func (mgr *AbiManager) SetAbi(account eos.AccountName, abi *eos.ABI) error {
	ctx, cancel := context.WithTimeout(mgr.ctx, time.Millisecond*500)
	defer cancel()
	return mgr.cache.Set(ctx, string(account), *abi, time.Hour)
}

// Get an ABI from the cache, on cache miss it is fetched from the
// API, gets cached and then returned to the user
func (mgr *AbiManager) GetAbi(account eos.AccountName) (*eos.ABI, error) {
	var abi eos.ABI
	if err := mgr.cacheGet(account, &abi); err != nil {
		ctx, cancel := context.WithTimeout(mgr.ctx, time.Second)
		defer cancel()
		resp, err := mgr.api.GetABI(ctx, account)
		if err != nil {
			return nil, fmt.Errorf("api: %s", err)
		}
		abi = resp.ABI

		err = mgr.SetAbi(account, &abi)
		if err != nil {
			return nil, fmt.Errorf("cache: %s", err)
		}
	}
	return &abi, nil
}

func (mgr *AbiManager) cacheGet(account eos.AccountName, value any) error {
	ctx, cancel := context.WithTimeout(mgr.ctx, time.Millisecond*500)
	defer cancel()
	return mgr.cache.Get(ctx, string(account), value)
}
