package config

import (
	"time"

	"github.com/eosswedenorg/thalos/internal/log"
	"github.com/eosswedenorg/thalos/internal/types"
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

type AbiCache struct {
	ApiTimeout time.Duration `yaml:"api_timeout" mapstructure:"api_timeout"`
}

type ShipConfig struct {
	Url                  string          `yaml:"url" mapstructure:"url"`
	IrreversibleOnly     bool            `yaml:"irreversible_only" mapstructure:"irreversible_only"`
	MaxMessagesInFlight  uint32          `yaml:"max_messages_in_flight" mapstructure:"max_messages_in_flight"`
	StartBlockNum        uint32          `yaml:"start_block_num" mapstructure:"start_block_num"`
	EndBlockNum          uint32          `yaml:"end_block_num" mapstructure:"end_block_num"`
	Chain                string          `yaml:"chain" mapstructure:"chain"`
	Blacklist            types.Blacklist `yaml:"blacklist" mapstructure:"blacklist"`
	BlacklistIsWhitelist bool            `yaml:"blacklist_is_whitelist" mapstructure:"blacklist_is_whitelist"`
}

type Config struct {
	Name string     `yaml:"name" mapstructure:"name"`
	Ship ShipConfig `yaml:"ship" mapstructure:"ship"`
	Api  string     `yaml:"api" mapstructure:"api"`

	Log log.Config `yaml:"log" mapstructure:"log"`

	Redis        RedisConfig `yaml:"redis" mapstructure:"redis"`
	MessageCodec string      `yaml:"message_codec" mapstructure:"message_codec"`

	AbiCache AbiCache `yaml:"abi_cache" mapstructure:"abi_cache"`

	Telegram TelegramConfig `yaml:"telegram" mapstructure:"telegram"`
}
