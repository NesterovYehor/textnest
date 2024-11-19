package jsonlog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
)

type Logger struct {
	logger *slog.Logger
}

// New creates a new logger instance with JSON output and a minimum logging level.
func New(out io.Writer, minLevel slog.Level) *Logger {
	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: minLevel})
	return &Logger{
		logger: slog.New(handler),
	}
}

// PrintInfo logs an informational message with additional properties.
func (l *Logger) PrintInfo(ctx context.Context, message string, properties map[string]string) {
	l.log(ctx, slog.LevelInfo, message, properties)
}

// PrintError logs an error message with additional properties.
func (l *Logger) PrintError(ctx context.Context, err error, properties map[string]string) {
	l.log(ctx, slog.LevelError, err.Error(), properties)
}

// PrintDebug logs a debug-level message with additional properties.
func (l *Logger) PrintDebug(ctx context.Context, message string, properties map[string]string) {
	l.log(ctx, slog.LevelDebug, message, properties)
}

// PrintFatal logs a fatal error message, adds a trace, and exits the application.
func (l *Logger) PrintFatal(ctx context.Context, err error, properties map[string]string) {
	l.log(ctx, slog.LevelError, err.Error(), properties)
	os.Exit(1)
}

// log handles logging messages at the specified level with optional properties.
func (l *Logger) log(ctx context.Context, level slog.Level, message string, properties map[string]string) {
	attrs := make([]slog.Attr, 0, len(properties)+1)

	// Add properties as structured attributes
	for k, v := range properties {
		attrs = append(attrs, slog.String(k, v))
	}

	// Add trace for error and fatal levels
	if level >= slog.LevelError {
		attrs = append(attrs, slog.String("trace", string(debug.Stack())))
	}

	// Log the message
	l.logger.Log(ctx, level, message, attrs)
}
