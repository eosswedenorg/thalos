package main

import (
	"encoding/json"
	"fmt"
	"time"

	"eosio-ship-trace-reader/internal/abi_cache"
	"eosio-ship-trace-reader/internal/redis"

	eos "github.com/eoscanada/eos-go"
	redis_cache "github.com/go-redis/cache/v8"
)

var abiCache *abi_cache.Cache

func InitAbiCache(id string) {
	// Init abi cache
	abiCache = abi_cache.New("ship.cache."+id+".abi", &redis_cache.Options{
		Redis: redis.Client(),
		// Cache 10k keys for 10 minutes.
		LocalCache: redis_cache.NewTinyLFU(10000, 10*time.Minute),
	})
}

func GetAbi(account eos.AccountName) (*eos.ABI, error) {
	key := string(account)

	abi, err := abiCache.Get(key)
	if err != nil {
		resp, err := eosClient.GetABI(eosClientCtx, account)
		if err != nil {
			return nil, fmt.Errorf("api: %s", err)
		}
		abi = &resp.ABI

		err = abiCache.Set(key, abi, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("cache: %s", err)
		}
	}
	return abi, nil
}

func DecodeAction(abi *eos.ABI, data []byte, actionName eos.ActionName) (interface{}, error) {
	var v interface{}

	bytes, err := abi.DecodeAction(data, actionName)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(bytes, &v)
	return v, err
}
