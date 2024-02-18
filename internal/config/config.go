package config

import (
	"time"

	"github.com/eosswedenorg/thalos/internal/log"

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
	Id      string `yaml:"id" mapstructure:"id"`
	Channel int64  `yaml:"channel" mapstructure:"channel"`
}

type ShipConfig struct {
	Url                 string `yaml:"url" mapstructure:"url"`
	IrreversibleOnly    bool   `yaml:"irreversible_only" mapstructure:"irreversible_only"`
	MaxMessagesInFlight uint32 `yaml:"max_messages_in_flight" mapstructure:"max_messages_in_flight"`
	StartBlockNum       uint32 `yaml:"start_block_num" mapstructure:"start_block_num"`
	EndBlockNum         uint32 `yaml:"end_block_num" mapstructure:"end_block_num"`
	Chain               string `yaml:"chain" mapstructure:"chain"`
}

type Config struct {
	Name string     `yaml:"name" mapstructure:"name"`
	Ship ShipConfig `yaml:"ship" mapstructure:"ship"`
	Api  string     `yaml:"api" mapstructure:"api"`

	Log log.Config `yaml:"log" mapstructure:"log"`

	Redis        RedisConfig `yaml:"redis" mapstructure:"redis"`
	MessageCodec string      `yaml:"message_codec" mapstructure:"message_codec"`

	Telegram TelegramConfig `yaml:"telegram" mapstructure:"telegram"`
}

// Create a new Config object with default values
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
