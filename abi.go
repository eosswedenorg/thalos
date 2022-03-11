
package main

import (
    "time"
    "errors"
    "fmt"
    "encoding/json"
    eos "github.com/eoscanada/eos-go"
    redis_cache "github.com/go-redis/cache/v8"
    "eosio-ship-trace-reader/abi_cache"
)

var abiCache *abi_cache.Cache

func InitAbiCache(id string) {
    // Init abi cache
    abiCache = abi_cache.New("ship.cache." + id + ".abi", &redis_cache.Options{
        Redis: rdb,
         // Cache 10k keys for 10 minutes.
        LocalCache: redis_cache.NewTinyLFU(10000, 10 * time.Minute),
    })
}

func GetAbi(account eos.AccountName) (*eos.ABI, error) {

    key := string(account)

    abi, err := abiCache.Get(key)
    if err != nil {
        resp, err := eosClient.GetABI(eosClientCtx, account)
        if err != nil {
            return nil, errors.New(fmt.Sprintf("api: %s", err))
        }
        abi = &resp.ABI

        err = abiCache.Set(key, abi, time.Hour)
        if err != nil {
            return nil, errors.New(fmt.Sprintf("cache: %s", err))
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
