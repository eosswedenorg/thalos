package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

// Rotating file represents a file that can be rotated when either the file
// becomes to large or to old, whatever comes first
type RotatingFile struct {
	fd      *os.File
	size    int64
	maxSize int64
	ts      time.Time
	maxAge  time.Duration
	format  string
}

func open(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o666)
}

// Open a new rotating file.
func NewRotatingFile(filename string, opts ...RotatingFileOption) (*RotatingFile, error) {
	if err := os.MkdirAll(path.Dir(filename), 0o766); err != nil && !os.IsExist(err) {
		return nil, err
	}

	fd, err := open(filename)
	if err != nil {
		return nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	file := &RotatingFile{
		fd:     fd,
		size:   stat.Size(),
		ts:     time.Now(),
		format: "2006-01-02_150405",
	}

	for _, opt := range opts {
		opt(file)
	}

	return file, nil
}

// Open a new rotating file using a config struct.
func NewRotatingFileFromConfig(config Config, suffix string) (*RotatingFile, error) {
	if len(suffix) > 0 {
		suffix = "_" + suffix
	}

	filename := config.GetFilePath() + suffix + ".log"

	return NewRotatingFile(filename, func(f *RotatingFile) {
		f.maxAge = config.MaxTime
		f.maxSize = int64(config.MaxFileSize)
	})
}

func (w *RotatingFile) newFilename(name string) string {
	ext := path.Ext(name)
	if len(ext) > 0 {
		name = name[:len(name)-len(ext)]
	}
	return fmt.Sprintf("%s-%s%s", name, time.Now().Format(w.format), ext)
}

// Get the filename
func (w RotatingFile) GetFilename() string {
	return path.Base(w.fd.Name())
}

// Rotate the file.
func (w *RotatingFile) Rotate() error {
	dst, err := os.OpenFile(w.newFilename(w.fd.Name()), os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Seek to the beginning of file
	if _, err = w.fd.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// And copy the contents to the new file.
	if _, err = io.Copy(dst, w.fd); err != nil {
		return err
	}

	// Then truncate the log.
	if err = w.fd.Truncate(0); err != nil {
		return err
	}

	w.size = 0
	w.ts = time.Now()

	return nil
}

// Implement io.Writer interface
func (w *RotatingFile) Write(p []byte) (int, error) {
	n, err := w.fd.Write(p)
	if err != nil {
		return n, err
	}

	w.size += int64(n)

	// Check if we should rotate
	if w.size >= w.maxSize || time.Since(w.ts) >= w.maxAge {
		if err := w.Rotate(); err != nil {
			return n, err
		}
	}

	return n, nil
}

// Implement io.Closer interface
func (w *RotatingFile) Close() error {
	err := w.fd.Close()
	w.fd = nil
	return err
}
