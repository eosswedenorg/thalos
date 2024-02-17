package config

import (
	"path"

	"github.com/spf13/pflag"
)

// Read cli flag values into the config
func (cfg *Config) ReadCliFlags(flags *pflag.FlagSet) error {
	logFile, _ := flags.GetString("log")
	if len(logFile) > 0 {
		cfg.Log.Directory = path.Dir(logFile)
		cfg.Log.Filename = path.Base(logFile)
	}
	return nil
}
