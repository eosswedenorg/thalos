package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

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

func NewRotatingFile(filename string, maxSize int64, maxAge time.Duration) (*RotatingFile, error) {
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

	return &RotatingFile{
		fd:      fd,
		size:    stat.Size(),
		maxSize: maxSize,
		ts:      time.Now(),
		maxAge:  maxAge,
		format:  "2006-01-02_150405",
	}, nil
}

func NewRotatingFileFromConfig(config Config) (*RotatingFile, error) {
	return NewRotatingFile(config.GetFilePath(), int64(config.MaxFileSize), config.MaxTime)
}

func (w *RotatingFile) newFilename(name string) string {
	ext := path.Ext(name)
	if len(ext) > 0 {
		name = name[:len(name)-len(ext)]
	}
	return fmt.Sprintf("%s-%s%s", name, time.Now().Format(w.format), ext)
}

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
