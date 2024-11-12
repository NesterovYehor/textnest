package logger

import (
	"log/slog"
	"os"
)

// Logger interface to be used by services
type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
}

// Config holds configuration options for the logger.
type Config struct {
	Level  slog.Level // Log level (e.g., Info, Warn, Error)
	Format string     // Format of the log output ("text" or "json")
	Output *os.File   // Log output destination (e.g., os.Stdout, os.Stderr)
}

// NewLogger creates a new logger instance with the given configuration.
func NewLogger(cfg Config) Logger {
	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(cfg.Output, &slog.HandlerOptions{Level: cfg.Level})
	} else {
		handler = slog.NewTextHandler(cfg.Output, &slog.HandlerOptions{Level: cfg.Level})
	}

	return slog.New(handler) // Return a new logger instance for the service
}

// Helper logging functions for simplicity
func Info(logger Logger, msg string, keyvals ...interface{}) {
	logger.Info(msg, keyvals...)
}

func Error(logger Logger, msg string, keyvals ...interface{}) {
	logger.Error(msg, keyvals...)
}

func Debug(logger Logger, msg string, keyvals ...interface{}) {
	logger.Debug(msg, keyvals...)
}
