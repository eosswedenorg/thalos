package config

import (
	"gopkg.in/yaml.v3"
)

func (ship *ShipConfig) UnmarshalYAML(value *yaml.Node) error {
	var err error

	if value.Kind == yaml.ScalarNode {
		ship.Url = value.Value
	} else {
		type ShipConfigRaw ShipConfig
		raw := ShipConfigRaw(*ship)
		if err = value.Decode(&raw); err == nil {
			*ship = ShipConfig(raw)
		}
	}

	return err
}

// Read YAML config data
func (cfg *Config) ReadYAML(data []byte) error {
	return yaml.Unmarshal(data, cfg)
}
