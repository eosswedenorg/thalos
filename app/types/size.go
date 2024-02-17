package types

import (
	"github.com/docker/go-units"
	"gopkg.in/yaml.v3"
)

// Size is an alias of int64 that can handle sizes represented
// in human readable strings like "200mb", "20 GB" etc.

// The value is in bytes.
type Size int64

// Parse a string into number of bytes stored in a int64
func (s *Size) Parse(value string) error {
	// Empty strings are not an error, they represents zero bytes.
	if len(value) < 1 {
		*s = 0
		return nil
	}

	v, err := units.FromHumanSize(value)
	if err != nil {
		return err
	}
	*s = Size(v)
	return nil
}

func (s Size) String() string {
	return units.HumanSize(float64(s))
}

func (s *Size) UnmarshalYAML(value *yaml.Node) error {
	return s.Parse(value.Value)
}

func (s *Size) UnmarshalText(text []byte) error {
	return s.Parse(string(text))
}
