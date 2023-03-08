package config

import (
	"encoding/json"
	"io/ioutil"

	"eosio-ship-trace-reader/app/service/redis"
	"eosio-ship-trace-reader/app/service/telegram"
)

const NULL_BLOCK_NUMBER uint32 = 0xffffffff

type Config struct {
	Name    string `json:"name"`
	ShipApi string `json:"ship_api"`
	Api     string `json:"api"`

	IrreversibleOnly    bool   `json:"irreversible_only"`
	MaxMessagesInFlight uint32 `json:"max_messages_in_flight"`
	StartBlockNum       uint32 `json:"start_block_num"`
	EndBlockNum         uint32 `json:"end_block_num"`

	Redis    redis.Config    `json:"redis"`
	Telegram telegram.Config `json:"telegram"`
}

func Parse(data []byte) (*Config, error) {
	cfg := Config{
		StartBlockNum:       NULL_BLOCK_NUMBER,
		EndBlockNum:         NULL_BLOCK_NUMBER,
		MaxMessagesInFlight: 10,
		IrreversibleOnly:    false,
		Redis:               redis.DefaultConfig,
	}

	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

func Load(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Parse(bytes)
}
