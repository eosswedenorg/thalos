package log

import "time"

type RotatingFileOption func(*RotatingFile)

func WithTimestampFormat(value string) RotatingFileOption {
	return func(f *RotatingFile) {
		if len(value) > 0 {
			f.format = value
		}
	}
}

func WithMaxSize(value int64) RotatingFileOption {
	return func(f *RotatingFile) {
		f.maxSize = value
	}
}

func WithMaxAge(value time.Duration) RotatingFileOption {
	return func(f *RotatingFile) {
		f.maxAge = value
	}
}
