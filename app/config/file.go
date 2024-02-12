package config

import (
	"os"
)

// Read values from file
func (cfg *Config) ReadFile(filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return cfg.ReadYAML(bytes)
}
