package audit

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileLogger wraps Logger and writes audit events to a file.
type FileLogger struct {
	*Logger
	file *os.File
}

// NewFileLogger opens (or creates) the given file path for appending
// and returns a FileLogger backed by it.
func NewFileLogger(path string) (*FileLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file %q: %w", path, err)
	}

	return &FileLogger{
		Logger: NewLogger(f),
		file:   f,
	}, nil
}

// Close flushes and closes the underlying log file.
func (fl *FileLogger) Close() error {
	if fl.file == nil {
		return nil
	}
	if err := fl.file.Sync(); err != nil {
		return fmt.Errorf("audit: sync log file: %w", err)
	}
	return fl.file.Close()
}
