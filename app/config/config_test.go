package config

import (
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/app/log"
	"github.com/stretchr/testify/require"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

func TestNew(t *testing.T) {
	expected := Config{
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
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Prefix:   "ship",
		},
	}

	require.Equal(t, expected, New())
}

func TestRead(t *testing.T) {
	expected := Config{
		Name:         "ship-reader-1",
		Api:          "http://127.0.0.1:8080",
		MessageCodec: "mojibake",
		Log: log.Config{
			Filename:    "some_file.log",
			Directory:   "/path/to/whatever",
			MaxFileSize: 200,
			MaxTime:     30 * time.Minute,
		},
		Ship: ShipConfig{
			Url:                 "127.0.0.1:8089",
			StartBlockNum:       23671836,
			EndBlockNum:         23872222,
			IrreversibleOnly:    true,
			MaxMessagesInFlight: 1337,
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

	cfg := Config{}
	err := cfg.Read([]byte(`
name: "ship-reader-1"
api: "http://127.0.0.1:8080"
message_codec: "mojibake"
log:
  filename: some_file.log
  directory: /path/to/whatever
  maxtime: 30m
  maxfilesize: 200b
ship:
  url: "127.0.0.1:8089"
  irreversible_only: true
  max_messages_in_flight: 1337
  start_block_num: 23671836
  end_block_num: 23872222
telegram:
  id: "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw"
  channel: -123456789
redis:
  addr: "localhost:6379"
  user: "myuser"
  password: "passwd"
  db: 4
  prefix: "some::ship"
`))

	require.NoError(t, err)
	require.Equal(t, expected, cfg)
}

func TestReadShorthandShipUrl(t *testing.T) {
	expected := Config{
		Name:         "ship-reader-1",
		Api:          "http://127.0.0.1:8080",
		MessageCodec: "json",
		Log: log.Config{
			MaxFileSize: 10 * 1000 * 1000, // 10 mb
			MaxTime:     time.Hour * 24,
		},
		Ship: ShipConfig{
			Url:                 "127.0.0.1:8089",
			StartBlockNum:       shipclient.NULL_BLOCK_NUMBER,
			EndBlockNum:         shipclient.NULL_BLOCK_NUMBER,
			MaxMessagesInFlight: 10,
			IrreversibleOnly:    false,
		},
		Telegram: TelegramConfig{
			Id:      "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw",
			Channel: -123456789,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "passwd",
			DB:       4,
			Prefix:   "some::ship",
		},
	}

	cfg := New()

	err := cfg.Read([]byte(`
name: "ship-reader-1"
api: "http://127.0.0.1:8080"
ship: "127.0.0.1:8089"
telegram:
  id: "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw"
  channel: -123456789
redis:
  addr: "localhost:6379"
  password: "passwd"
  db: 4
  prefix: "some::ship"
`))

	require.NoError(t, err)
	require.Equal(t, expected, cfg)
}
