package log

import (
	"testing"
)

func TestConfig_GetDirectory(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		want      string
	}{
		{"empty", "", "."},
		{"root", "/", "/"},
		{"one", "dir", "dir"},
		{"path", "/path/to/some/directory", "/path/to/some/directory"},
		{"relative", "relative/directory", "relative/directory"},
		{"backtrace", "/path/./to/some/../directory", "/path/to/directory"},
		{"multislash", "//path/to///directory//", "/path/to/directory"},
		{"everything", "path/to/..//./from/directory//", "path/from/directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Directory: tt.directory,
			}
			if got := c.GetDirectory(); got != tt.want {
				t.Errorf("Config.GetDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{"empty", "", "."},
		{"name", "some_file.txt", "some_file.txt"},
		{"path", "/path/to/my.log", "my.log"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Filename: tt.filename,
			}
			if got := c.GetFilename(); got != tt.want {
				t.Errorf("Config.GetFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetFilePath(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		directory string
		want      string
	}{
		{"empty", "", "", "."},
		{"directory", "", "dir", "dir"},
		{"filename", "filename", "", "filename"},
		{"both", "filename", "dir", "dir/filename"},
		{"root", "filename", "/", "/filename"},
		{"abs", "filename", "/path/to/logs", "/path/to/logs/filename"},
		{"relative", "filename", "/srv/../log", "/log/filename"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Filename:  tt.filename,
				Directory: tt.directory,
			}
			if got := c.GetFilePath(); got != tt.want {
				t.Errorf("Config.GetFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
