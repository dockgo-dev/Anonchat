package lib

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	logDirectory = "logs"
	logLayout    = "2006-01-02"
)

// NewLogger constructs a slog.Logger that writes to both stdout and a daily rotated file.
// The returned cleanup function must be called to close underlying resources when done.
func NewLogger(name string) (*slog.Logger, func() error, error) {
	if name == "" {
		name = "app"
	}

	if err := os.MkdirAll(logDirectory, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create log directory: %w", err)
	}

	filename := filepath.Join(logDirectory, fmt.Sprintf("%s-%s.log", name, time.Now().Format(logLayout)))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	logger := slog.New(handler).With("component", name)

	cleanup := func() error {
		return file.Close()
	}

	return logger, cleanup, nil
}
