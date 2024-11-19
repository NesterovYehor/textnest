package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	log *slog.Logger
}

// NewLogger creates a new logger instance with the given configuration.
func NewLogger(w io.Writer, level slog.Level) *Logger {
	opts := slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(w, &opts))
	slog.SetDefault(logger) // Set this as the global default logger if needed.

	return &Logger{
		log: logger,
	}
}

// DefaultLogger initializes a default logger writing to stdout at the INFO level.
func DefaultLogger() *Logger {
	return NewLogger(os.Stdout, slog.LevelInfo)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.log.Info(msg, keysAndValues...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, keysAndValues ...any) {
	l.log.Error(msg, keysAndValues...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, keysAndValues ...any) {
	l.log.Warn(msg, keysAndValues...)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.log.Debug(msg, keysAndValues...)
}
