package config

import (
	"path"

	"github.com/urfave/cli/v2"
)

// Read cli flag values into the config
func (cfg *Config) ReadCliFlags(ctx *cli.Context) error {
	logFile := ctx.Path("log")
	if len(logFile) > 0 {
		cfg.Log.Directory = path.Dir(logFile)
		cfg.Log.Filename = path.Base(logFile)
	}

	return nil
}
