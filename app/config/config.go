package config

import (
	"time"

	"github.com/eosswedenorg/thalos/app/log"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Prefix   string `yaml:"prefix"`
}

type TelegramConfig struct {
	Id      string `yaml:"id"`
	Channel int64  `yaml:"channel"`
}

type ShipConfig struct {
	Url                 string `yaml:"url"`
	IrreversibleOnly    bool   `yaml:"irreversible_only"`
	MaxMessagesInFlight uint32 `yaml:"max_messages_in_flight"`
	StartBlockNum       uint32 `yaml:"start_block_num"`
	EndBlockNum         uint32 `yaml:"end_block_num"`
	Chain               string `yaml:"chain"`
}

type Config struct {
	Name string     `yaml:"name"`
	Ship ShipConfig `yaml:"ship"`
	Api  string     `yaml:"api"`

	Log log.Config `yaml:"log"`

	Redis        RedisConfig `yaml:"redis"`
	MessageCodec string      `yaml:"message_codec"`

	Telegram TelegramConfig `yaml:"telegram"`
}

func New() Config {
	return Config{
		MessageCodec: "json",
		Log: log.Config{
			MaxFileSize: 10 * 1000 * 1000, // 10 mb
			MaxTime:     time.Hour * 24,
		},
		Ship: ShipConfig{
			StartBlockNum:       shipclient.NULL_BLOCK_NUMBER,
			EndBlockNum:         shipclient.NULL_BLOCK_NUMBER,
			MaxMessagesInFlight: 10,
			IrreversibleOnly:    false,
		},
		Redis: RedisConfig{
			Addr:   "localhost:6379",
			Prefix: "ship",
		},
	}
}
