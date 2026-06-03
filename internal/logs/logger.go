package logs

import (
	"io"
	"log/slog"
	"os"
)

// Setup sets up the global slog logger to write JSON to both stdout and a file.
func Setup(logPath string) (*os.File, error) {
	// Create log file if it doesn't exist, open for appending
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Use MultiWriter to write to both stdout and the file
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Create JSON handler
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Create and set global logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return file, nil
}
