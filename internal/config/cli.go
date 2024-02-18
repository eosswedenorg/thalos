package config

import (
	"path"

	"github.com/spf13/pflag"
)

func GetFlags() *pflag.FlagSet {
	flags := pflag.FlagSet{}
	flags.StringP("config", "c", "./config.yml", "Config file to read")
	flags.StringP("level", "L", "info", "Log level to use")
	flags.StringP("log", "l", "", "Path to log file (default: print to stdout/stderr)")
	flags.StringP("pid", "p", "", "Where to write process id")
	flags.BoolP("no-state-cache", "n", false, "Force the application to take start block from config/api")

	flags.Int("start-block", 0, "Start to stream from this block (default: config value, cache, head from api)")
	flags.Int("end-block", 0, "Stop streaming when this block is reached")

	return &flags
}

func overrideCliFlags(cfg *Config, flags *pflag.FlagSet) error {
	logFile, _ := flags.GetString("log")
	if len(logFile) > 0 {
		cfg.Log.Directory = path.Dir(logFile)
		cfg.Log.Filename = path.Base(logFile)
	}
	return nil
}
