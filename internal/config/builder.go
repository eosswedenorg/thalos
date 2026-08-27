package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/eosswedenorg/thalos/internal/types"
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

			"cache.storage": "cache",

			// AbiCache
			"abi_cache.api_timeout": "abi-cache-api-timeout",

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
			"ship.blacklist":              "blacklist",
			"ship.blacklist_is_whitelist": "blacklist-is-whitelist",
			"ship.table_deltas":           "table-deltas",
			"ship.table_delta_whitelist":  "table-delta-whitelist",
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
		func(f reflect.Type, t reflect.Type, in interface{}) (interface{}, error) {
			if t == reflect.TypeOf(types.Blacklist{}) {
				return decodeIntoBlacklist(in)
			}
			return in, nil
		},
	)

	err := v.Unmarshal(&conf, viper.DecodeHook(decoders))
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Decode a generic structure into types.Blacklist
func decodeIntoBlacklist(in any) (*types.Blacklist, error) {
	switch v := in.(type) {
	// Standard map structure.
	case map[string]any:
		return blacklistParseMap(v)

	// slice of "contract:action" pairs. Usually from CLI
	case []string:
		return blacklistParseSlice(v)

	// Sometimes we have a slice of interfaces.
	// Need to convert it to a slice of strings.
	case []any:
		sv, err := toStringSlice(v)
		if err != nil {
			return nil, err
		}
		return blacklistParseSlice(sv)
	}

	return nil, fmt.Errorf("Must be a string slice")
}

// Blacklist map parser
func blacklistParseMap(in map[string]any) (*types.Blacklist, error) {
	list := &types.Blacklist{}
	for contract, value := range in {
		if err := blacklistParseValue(list, contract, value); err != nil {
			return nil, err
		}
	}
	return list, nil
}

func blacklistParseValue(list *types.Blacklist, contract string, in any) error {
	switch value := in.(type) {
	case map[string]any:
		for nested, nestedValue := range value {
			if err := blacklistParseValue(list, contract+"."+nested, nestedValue); err != nil {
				return err
			}
		}
		return nil
	case map[any]any:
		for nested, nestedValue := range value {
			nestedKey, ok := nested.(string)
			if !ok {
				return fmt.Errorf("Must be a string slice")
			}
			if err := blacklistParseValue(list, contract+"."+nestedKey, nestedValue); err != nil {
				return err
			}
		}
		return nil
	}

	values, err := toStringSlice(in)
	if err != nil {
		return err
	}

	for _, entry := range values {
		list.Add(contract, entry)
	}

	return nil
}

// Blacklist slice parser
func blacklistParseSlice(in []string) (*types.Blacklist, error) {
	list := &types.Blacklist{}
	for _, i := range in {
		var action string
		parts := strings.SplitN(i, ":", 2)

		if len(parts) < 2 {
			action = "*"
		} else {
			action = parts[1]
		}

		list.Add(parts[0], action)
	}
	return list, nil
}

func toStringSlice(in any) ([]string, error) {
	switch value := in.(type) {
	case string:
		return []string{value}, nil
	case []string:
		return value, nil
	case []any:
		out := make([]string, 0, len(value))
		for _, item := range value {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("Must be a string slice")
			}
			out = append(out, str)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("Must be a string slice")
	}
}
