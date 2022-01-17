
package main

import (
    "time"
    "encoding/json"
    eos "github.com/eoscanada/eos-go"
    "github.com/go-redis/redis/v8"
)

func GetAbi(account eos.AccountName) (*eos.ABI, error) {

    key := RedisKey("abi", string(account))

    data, err := RedisGet(key).Result()
    if err == redis.Nil {
        val, err := eosClient.GetABI(eosClientCtx, account)
        if err != nil {
            return nil, err
        }

        b, err := json.Marshal(val.ABI)
        data = string(b)

        err = RedisSet(key, data, time.Hour).Err()
        if err != nil {
            return nil, err
        }
    }

    abi := eos.ABI{}
    err = json.Unmarshal([]byte(data), &abi)
    if err != nil {
        return nil, err
    }
    return &abi, nil
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
