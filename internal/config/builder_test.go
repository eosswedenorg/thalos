package config

import (
	"bytes"
	"testing"
	"time"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
	"github.com/eosswedenorg/thalos/internal/log"
	"github.com/eosswedenorg/thalos/internal/types"
	"github.com/karlseguin/typed"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	expected := Config{
		Name:         "ship-reader-1",
		Api:          "http://127.0.0.1:8080",
		MessageCodec: "mojibake",
		Log: log.Config{
			Filename:            "some_file.log",
			Directory:           "/path/to/whatever",
			MaxFileSize:         200,
			MaxTime:             30 * time.Minute,
			FileTimestampFormat: "20060102@150405",
		},
		Cache: Cache{
			Storage: "memcached",
			Options: typed.Typed{
				"ttl":             "300m",
				"size":            400,
				"super_fast_mode": true,
			},
		},
		AbiCache: AbiCache{
			ApiTimeout: time.Minute * 300,
		},
		Ship: ShipConfig{
			Url:                 "127.0.0.1:8089",
			StartBlockNum:       23671836,
			EndBlockNum:         23872222,
			IrreversibleOnly:    true,
			MaxMessagesInFlight: 1337,
			Blacklist: *types.NewBlacklist(map[string][]string{
				"eosio":    {"noop"},
				"contract": {"skip1", "skip2"},
			}),
			BlacklistIsWhitelist: true,
		},
		Telegram: TelegramConfig{
			Id:      "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw",
			Channel: -123456789,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			User:     "myuser",
			Password: "passwd",
			DB:       4,
			Prefix:   "some::ship",
		},
	}

	builder := NewBuilder()
	builder.SetSource(bytes.NewBuffer([]byte(`
name: "ship-reader-1"
api: "http://127.0.0.1:8080"
message_codec: "mojibake"
cache:
  storage: memcached
  options:
    ttl: 300m
    size: 400
    super_fast_mode: true
abi_cache:
  api_timeout: 300m
log:
  filename: some_file.log
  directory: /path/to/whatever
  maxtime: 30m
  maxfilesize: 200b
  file_timestamp_format: 20060102@150405
ship:
  url: "127.0.0.1:8089"
  irreversible_only: true
  max_messages_in_flight: 1337
  start_block_num: 23671836
  end_block_num: 23872222
  blacklist:
    eosio: noop
    contract:
      - skip1
      - skip2
  blacklist_is_whitelist: true
telegram:
  id: "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw"
  channel: -123456789
redis:
  addr: "localhost:6379"
  user: "myuser"
  password: "passwd"
  db: 4
  prefix: "some::ship"
`)))

	cfg, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, &expected, cfg)
}

func TestBuilder_WithDefaultConfig(t *testing.T) {
	expected := Config{
		MessageCodec: "json",
		Log: log.Config{
			MaxFileSize:         10 * 1000 * 1000,
			MaxTime:             time.Hour * 24,
			FileTimestampFormat: "2006-01-02_150405",
		},
		Cache: Cache{
			Storage: "redis",
		},
		AbiCache: AbiCache{
			ApiTimeout: time.Second,
		},
		Ship: ShipConfig{
			Url:                 "ws://127.0.0.1:8080",
			StartBlockNum:       shipclient.NULL_BLOCK_NUMBER,
			EndBlockNum:         shipclient.NULL_BLOCK_NUMBER,
			MaxMessagesInFlight: 10,
			EnableTableDeltas:   true,
		},
		Redis: RedisConfig{
			Addr:   "127.0.0.1:6379",
			Prefix: "ship",
		},
	}

	cfg, err := NewBuilder().
		SetSource(bytes.NewReader([]byte(``))).
		SetFlags(GetFlags()).
		Build()

	require.NoError(t, err)
	require.Equal(t, &expected, cfg)
}

func TestBuilder_NilSource(t *testing.T) {
	cfg, err := NewBuilder().Build()
	require.Nil(t, cfg)
	require.EqualError(t, err, "Config not set")
}

func TestBuilder_Flags(t *testing.T) {
	flags := GetFlags()

	require.NoError(t, flags.Set("url", "https://myapi"))
	require.NoError(t, flags.Set("codec", "binary"))
	require.NoError(t, flags.Set("redis-addr", "154.223.38.15:6380"))
	require.NoError(t, flags.Set("redis-user", "myuser"))
	require.NoError(t, flags.Set("redis-password", "secret123"))
	require.NoError(t, flags.Set("redis-db", "3"))
	require.NoError(t, flags.Set("redis-prefix", "custom-prefix"))
	require.NoError(t, flags.Set("cache", "memcached"))
	require.NoError(t, flags.Set("abi-cache-api-timeout", "16h"))
	require.NoError(t, flags.Set("telegram-id", "72983126312982618"))
	require.NoError(t, flags.Set("telegram-channel", "-293492332"))
	require.NoError(t, flags.Set("log-max-filesize", "25mb"))
	require.NoError(t, flags.Set("log-max-time", "10m"))
	require.NoError(t, flags.Set("log-file-timestamp", "0102-15:04:05"))
	require.NoError(t, flags.Set("ship-url", "ws://myship.com:7823"))
	require.NoError(t, flags.Set("start-block", "7327833"))
	require.NoError(t, flags.Set("end-block", "329408392"))
	require.NoError(t, flags.Set("irreversible-only", "true"))
	require.NoError(t, flags.Set("max-msg-in-flight", "98"))
	require.NoError(t, flags.Set("chain", "wax"))
	require.NoError(t, flags.Set("blacklist", "contract:action1,contract:action2,contract2:action1"))
	require.NoError(t, flags.Set("blacklist-is-whitelist", "true"))
	require.NoError(t, flags.Set("table-deltas", "false"))

	cfg, err := NewBuilder().
		SetSource(bytes.NewReader([]byte(``))).
		SetFlags(flags).
		Build()

	expected := Config{
		Api:          "https://myapi",
		MessageCodec: "binary",
		Log: log.Config{
			MaxFileSize:         25 * 1000 * 1000, // 25 mb
			MaxTime:             time.Minute * 10,
			FileTimestampFormat: "0102-15:04:05",
		},
		Cache: Cache{
			Storage: "memcached",
		},
		AbiCache: AbiCache{
			ApiTimeout: time.Hour * 16,
		},
		Ship: ShipConfig{
			Url:                 "ws://myship.com:7823",
			StartBlockNum:       7327833,
			EndBlockNum:         329408392,
			MaxMessagesInFlight: 98,
			IrreversibleOnly:    true,
			Chain:               "wax",
			Blacklist: *types.NewBlacklist(map[string][]string{
				"contract":  {"action1", "action2"},
				"contract2": {"action1"},
			}),
			BlacklistIsWhitelist: true,
		},
		Telegram: TelegramConfig{
			Id:      "72983126312982618",
			Channel: -293492332,
		},
		Redis: RedisConfig{
			Addr:     "154.223.38.15:6380",
			User:     "myuser",
			Password: "secret123",
			DB:       3,
			Prefix:   "custom-prefix",
		},
	}

	require.NoError(t, err)
	require.Equal(t, &expected, cfg)
}

func TestBuilder_BlacklistSlice(t *testing.T) {
	expected := Config{
		Ship: ShipConfig{
			Blacklist: *types.NewBlacklist(map[string][]string{
				"contract":  {"action"},
				"contract2": {"action2"},
				"contract3": {"*"},
			}),
		},
	}

	builder := NewBuilder()
	builder.SetSource(bytes.NewBuffer([]byte(`
ship:
  blacklist:
    - "contract:action"
    - "contract2:action2"
    - contract3
`)))

	cfg, err := builder.Build()

	require.NoError(t, err)
	require.Equal(t, &expected, cfg)
}
