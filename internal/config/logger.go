package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/dbunt1tled/url-shortener/internal/lib/logger"
	"github.com/dbunt1tled/url-shortener/internal/lib/logger/handlers/slogpretty"
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type AppLogger struct {
	*slog.Logger

	level      slog.Level
	fileWriter *os.File
}

var (
	logInstance *AppLogger //nolint:gochecknoglobals // singleton
	lm          sync.Once  //nolint:gochecknoglobals // singleton
)

func LoadLogger(env string, debugLevel slog.Level) *AppLogger {
	lm.Do(func() {
		logInstance = &AppLogger{
			Logger: initLogger(env, debugLevel),
			level:  LevelDebug,
		}
	})
	return logInstance
}

func initLogger(env string, level slog.Level) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvDev:
		log = prettyLogHandler(env, level)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
		)
	}

	return log
}

func prettyLogHandler(env string, level slog.Level) *slog.Logger {

	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler.WithAttrs([]slog.Attr{
		slog.String("env", env),
	}))
}

// WithContext returns a logger with context values.
func (l *AppLogger) WithContext(ctx context.Context) *slog.Logger {
	// Extract values from context and add them to the logger
	// This is a placeholder implementation that can be expanded
	// to extract specific values from the context
	return l.Logger.With("context", "true")
}

// With returns a logger with the given attributes.
func (l *AppLogger) With(args ...any) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.With(args...),
		level:  l.level,
	}
	return newLogger
}

// WithGroup returns a logger with the given group.
func (l *AppLogger) WithGroup(name string) *AppLogger {
	newLogger := &AppLogger{
		Logger: l.Logger.WithGroup(name),
		level:  l.level,
	}
	return newLogger
}

// ErrorWithStack logs an error with its stack trace if available.
func (l *AppLogger) ErrorWithStack(msg string, err error) {
	attrs := logger.Error(err)
	args := make([]any, 0, len(attrs)*2) //nolint:mnd // dual volume
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	l.Logger.Error(msg, args...)
}

// Error logs an error message with structured data.
func (l *AppLogger) Error(msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2) //nolint:mnd // dual volume
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.Error(msg, args...)
}

// ErrorContext logs an error message with context and structured data.
func (l *AppLogger) ErrorContext(ctx context.Context, msg string, err error, keyvals ...any) {
	args := make([]any, 0, len(keyvals)+2) //nolint:mnd // dual volume
	args = append(args, "error", err.Error())
	args = append(args, keyvals...)
	l.Logger.ErrorContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and structured data.
func (l *AppLogger) WarnContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.WarnContext(ctx, msg, keyvals...)
}

// InfoContext logs an info message with context and structured data.
func (l *AppLogger) InfoContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.InfoContext(ctx, msg, keyvals...)
}

// DebugContext logs a debug message with context and structured data.
func (l *AppLogger) DebugContext(ctx context.Context, msg string, keyvals ...any) {
	l.Logger.DebugContext(ctx, msg, keyvals...)
}

// Fatalf logs a fatal error message and exits the program.
func (l *AppLogger) Fatalf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Errorf logs an error message.
func (l *AppLogger) Errorf(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// Warnf logs a warning message.
func (l *AppLogger) Warnf(msg string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// Infof logs an info message.
func (l *AppLogger) Infof(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// Debugf logs a debug message.
func (l *AppLogger) Debugf(msg string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}
