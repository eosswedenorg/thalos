package config

import (
	"errors"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// This is a simple module that encapsulate the creation
// of a config object and can override values from cli flags.

type Builder struct {
	in    io.Reader
	flags *pflag.FlagSet
	binds map[string]string
}

func NewBuilder() *Builder {
	return &Builder{
		binds: map[string]string{
			"api":           "url",
			"message_codec": "codec",

			// Redis
			"redis.addr":     "redis-addr",
			"redis.user":     "redis-user",
			"redis.password": "redis-password",
			"redis.db":       "redis-db",
			"redis.prefix":   "redis-prefix",

			// Telegram
			"telegram.id":      "telegram-id",
			"telegram.channel": "telegram-channel",

			// Log
			"log.maxfilesize":           "log-max-filesize",
			"log.maxtime":               "log-max-time",
			"log.file_timestamp_format": "log-file-timestamp",

			// Ship
			"ship.url":                    "ship-url",
			"ship.start_block_num":        "start-block",
			"ship.end_block_num":          "end-block",
			"ship.irreversible_only":      "irreversible-only",
			"ship.max_messages_in_flight": "max-msg-in-flight",
			"ship.chain":                  "chain",
		},
	}
}

// Set the config file to read
func (b *Builder) SetConfigFile(filename string) *Builder {
	file, _ := os.Open(filename)
	return b.SetSource(file)
}

// Set the source to read
func (b *Builder) SetSource(in io.Reader) *Builder {
	b.in = in
	return b
}

// Set all flags that the builder should use.
func (b *Builder) SetFlags(flags *pflag.FlagSet) *Builder {
	b.flags = flags
	return b
}

// Add a flag to the builder.
func (b *Builder) AddFlag(flag *pflag.Flag) *Builder {
	b.flags.AddFlag(flag)
	return b
}

// Build the config object from file, cli-flags
func (b *Builder) Build() (*Config, error) {
	if b.in == nil {
		return nil, errors.New("Config not set")
	}

	conf := Config{}

	v := viper.New()
	v.SetConfigType("yaml")

	if b.flags != nil {
		// bind flags in viper.
		for key, flagname := range b.binds {
			flag := b.flags.Lookup(flagname)
			if flag == nil {
				continue
			}

			if err := v.BindPFlag(key, flag); err != nil {
				return nil, err
			}
		}
	}

	// Read config and unmarshal
	if err := v.ReadConfig(b.in); err != nil {
		return nil, err
	}

	decoders := mapstructure.ComposeDecodeHookFunc(
		mapstructure.TextUnmarshallerHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)

	err := v.Unmarshal(&conf, viper.DecodeHook(decoders))
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
