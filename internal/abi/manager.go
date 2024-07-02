package abi

import (
	"context"
	"fmt"
	"time"

	"github.com/eosswedenorg/thalos/internal/cache"
	"github.com/shufflingpixels/antelope-go/api"
	"github.com/shufflingpixels/antelope-go/chain"
)

// AbiManager handles an ABI cache that fetches the ABI from an API on cache miss.
type AbiManager struct {
	cache *cache.Cache
	api   *api.Client
	ctx   context.Context
}

// Create a new ABI Manager
func NewAbiManager(cache *cache.Cache, api *api.Client) *AbiManager {
	return &AbiManager{
		cache: cache,
		api:   api,
		ctx:   context.Background(),
	}
}

// Set or update an ABI in the cache.
func (mgr *AbiManager) SetAbi(account chain.Name, abi *chain.Abi) error {
	ctx, cancel := context.WithTimeout(mgr.ctx, time.Millisecond*500)
	defer cancel()
	return mgr.cache.Set(ctx, account.String(), *abi, time.Hour)
}

// Get an ABI from the cache, on cache miss it is fetched from the
// API, gets cached and then returned to the user
func (mgr *AbiManager) GetAbi(account chain.Name) (*chain.Abi, error) {
	var abi chain.Abi
	if err := mgr.cacheGet(account, &abi); err != nil {
		ctx, cancel := context.WithTimeout(mgr.ctx, time.Second)
		defer cancel()
		resp, err := mgr.api.GetAbi(ctx, account.String())
		if err != nil {
			return nil, fmt.Errorf("api: %s", err)
		}
		abi = resp.Abi

		err = mgr.SetAbi(account, &abi)
		if err != nil {
			return nil, fmt.Errorf("cache: %s", err)
		}
	}
	return &abi, nil
}

func (mgr *AbiManager) cacheGet(account chain.Name, value any) error {
	ctx, cancel := context.WithTimeout(mgr.ctx, time.Millisecond*500)
	defer cancel()
	return mgr.cache.Get(ctx, account.String(), value)
}
