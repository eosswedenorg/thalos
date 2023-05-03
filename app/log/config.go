package log

import (
	"path"
	"time"

	"github.com/eosswedenorg/thalos/app/types"
)

type Config struct {
	Filename    string        `yaml:"filename"`
	Directory   string        `yaml:"directory"`
	MaxFileSize types.Size    `yaml:"maxfilesize"`
	MaxTime     time.Duration `yaml:"maxtime"`
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
