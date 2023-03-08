package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

func TestParse_Default(t *testing.T) {
	expected := Config{
		StartBlockNum:       shipclient.NULL_BLOCK_NUMBER,
		EndBlockNum:         shipclient.NULL_BLOCK_NUMBER,
		MaxMessagesInFlight: 10,
		IrreversibleOnly:    false,
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Prefix:   "ship",
		},
	}

	cfg, err := Parse([]byte(`{}`))
	require.NoError(t, err)
	require.Equal(t, cfg, &expected)
}

func TestParse(t *testing.T) {
	expected := Config{
		Name:                "ship-reader-1",
		Api:                 "http://127.0.0.1:8080",
		ShipApi:             "127.0.0.1:8089",
		StartBlockNum:       23671836,
		EndBlockNum:         23872222,
		IrreversibleOnly:    true,
		MaxMessagesInFlight: 1337,
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

	cfg, err := Parse([]byte(`{
		"name": "ship-reader-1",
		"api": "http://127.0.0.1:8080",
		"ship_api": "127.0.0.1:8089",
		"irreversible_only": true,
		"max_messages_in_flight": 1337,
		"start_block_num": 23671836,
		"end_block_num": 23872222,
		"telegram": {
			"id": "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw",
			"channel": -123456789
		},
		"redis": {
			"addr": "localhost:6379",
			"password": "passwd",
			"db": 4,
			"prefix": "some::ship"
		}
	}`))

	require.NoError(t, err)
	require.Equal(t, cfg, &expected)
}
