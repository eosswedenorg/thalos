package config

import (
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/internal/log"
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
