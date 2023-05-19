package log

import (
	"io"

	log "github.com/sirupsen/logrus"
)

type HookWriter struct {
	Writer    io.Writer
	LogLevels []log.Level
}

func (hook *HookWriter) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *HookWriter) Levels() []log.Level {
	return hook.LogLevels
}

func MakeStdHook(writer io.Writer) *HookWriter {
	return &HookWriter{
		Writer: writer,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
		},
	}
}

func MakeErrorHook(writer io.Writer) *HookWriter {
	return &HookWriter{
		Writer: writer,
		LogLevels: []log.Level{
			log.ErrorLevel,
			log.WarnLevel,
			log.FatalLevel,
			log.PanicLevel,
			log.TraceLevel,
		},
	}
}
