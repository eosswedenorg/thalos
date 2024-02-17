package config

import (
	"bytes"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Read values from file
func (cfg *Config) ReadFile(filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return cfg.Read(bytes)
}

func (cfg *Config) Read(in []byte) error {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer(in)); err != nil {
		return err
	}

	decoders := mapstructure.ComposeDecodeHookFunc(
		mapstructure.TextUnmarshallerHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		decodeShorthandShipConfig,
	)

	return v.Unmarshal(cfg, viper.DecodeHook(decoders))
}
