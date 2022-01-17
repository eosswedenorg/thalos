
package main

import (
    "io/ioutil"
    "encoding/json"
)

const NULL_BLOCK_NUMBER uint32 = 0xffffffff

type RedisConfig struct {
    Addr string `json:"addr"`
    Password string `json:"password"`
    DB int `json:db`
}

type Config struct {
    ShipApi string `json:"ship_api"`
    Api string `json:"api"`

    Redis RedisConfig `json:"redis"`

    IrreversibleOnly bool `json:"irreversible_only"`
    MaxMessagesInFlight uint32 `json:"max_messages_in_flight"`
    StartBlockNum uint32 `json:"start_block_num"`
    EndBlockNum uint32 `json:"end_block_num"`
}

func LoadConfig(filename string) (Config, error) {

    cfg := Config{
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
