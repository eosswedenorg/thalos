
package config

import (
    "io/ioutil"
    "encoding/json"
)

const NULL_BLOCK_NUMBER uint32 = 0xffffffff

type RedisConfig struct {
    Addr string `json:"addr"`
    Password string `json:"password"`
    DB int `json:db`
    CacheID string `json:"cache_id"`
}

type TelegramConfig struct {
    Id string `json:"id"`
    Channel int64 `json:"channel"`
}

type Config struct {
    Name string `json:"name"`
    ShipApi string `json:"ship_api"`
    Api string `json:"api"`
    Transport string `json:"transport"`

    Redis RedisConfig `json:"redis"`

    Telegram TelegramConfig `json:"telegram"`

    IrreversibleOnly bool `json:"irreversible_only"`
    MaxMessagesInFlight uint32 `json:"max_messages_in_flight"`
    StartBlockNum uint32 `json:"start_block_num"`
    EndBlockNum uint32 `json:"end_block_num"`
}

func Load(filename string) (Config, error) {

    cfg := Config{
        Transport: "redis-channel",
        StartBlockNum: NULL_BLOCK_NUMBER,
        EndBlockNum: NULL_BLOCK_NUMBER,
        MaxMessagesInFlight: 10,
        IrreversibleOnly: false,
        Redis: RedisConfig{
            Addr: "localhost:6379",
            Password: "",
            DB: 0,
        },
    }

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return cfg, err
    }

    err = json.Unmarshal(bytes, &cfg)
    if err != nil {
        return cfg, err
    }

    return cfg, nil
}
