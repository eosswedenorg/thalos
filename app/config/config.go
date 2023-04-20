package config

import (
	"encoding/json"
	"io/ioutil"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	CacheID  string `json:"cache_id"`
	Prefix   string `json:"prefix"`
}

type TelegramConfig struct {
	Id      string `json:"id"`
	Channel int64  `json:"channel"`
}

type ShipConfig struct {
	Url                 string `json:"url"`
	IrreversibleOnly    bool   `json:"irreversible_only"`
	MaxMessagesInFlight uint32 `json:"max_messages_in_flight"`
	StartBlockNum       uint32 `json:"start_block_num"`
	EndBlockNum         uint32 `json:"end_block_num"`
}

type Config struct {
	Name string     `json:"name"`
	Ship ShipConfig `json:"ship"`
	Api  string     `json:"api"`

	Redis RedisConfig `json:"redis"`

	Telegram TelegramConfig `json:"telegram"`
}

func Parse(data []byte) (*Config, error) {
	cfg := Config{
		Ship: ShipConfig{
			StartBlockNum:       shipclient.NULL_BLOCK_NUMBER,
			EndBlockNum:         shipclient.NULL_BLOCK_NUMBER,
			MaxMessagesInFlight: 10,
			IrreversibleOnly:    false,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Prefix:   "ship",
		},
	}

	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

func (ship *ShipConfig) UnmarshalJSON(data []byte) error {
	var err error

	if err = json.Unmarshal(data, &ship.Url); err != nil {
		//
		type ShipConfigRaw ShipConfig
		raw := ShipConfigRaw(*ship)
		if err = json.Unmarshal(data, &raw); err == nil {
			*ship = ShipConfig(raw)
		}
	}
	return err
}

func Load(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return Parse(bytes)
}
