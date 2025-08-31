package logger

import (
	"context"
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init initializes the default logger with configuration
func Init(level slog.Level, format string) {
	var handler slog.Handler

	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		})
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Get returns the default logger instance
func Get() *slog.Logger {
	if defaultLogger == nil {
		// Fallback to default if not initialized
		defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return defaultLogger
}

// Convenience methods that use the default logger
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Context-aware logging methods
func DebugContext(ctx context.Context, msg string, args ...any) {
	Get().DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Get().InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	Get().WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Get().ErrorContext(ctx, msg, args...)
}
