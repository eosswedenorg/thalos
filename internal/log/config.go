package log

import (
	"path"
	"time"

	"github.com/eosswedenorg/thalos/internal/types"
)

// Config represents configuration parameters for a log.
type Config struct {
	// Filename where the log is stored.
	Filename string `yaml:"filename" mapstructure:"filename"`

	// Directory where the log files are stored.
	Directory string `yaml:"directory" mapstructure:"directory"`

	// Maximum filesize, the log is rotated when this size is exceeded.
	MaxFileSize types.Size `yaml:"maxfilesize" mapstructure:"maxfilesize"`

	// Maximum lifetime of the file before it is rotated.
	MaxTime time.Duration `yaml:"maxtime" mapstructure:"maxtime"`
}

func (c Config) GetFilename() string {
	return path.Base(c.Filename)
}

func (c Config) GetDirectory() string {
	return path.Clean(c.Directory)
}

func (c Config) GetFilePath() string {
	return path.Join(c.GetDirectory(), c.GetFilename())
}
