package config

import (
	"go_first/internal/lib/logger/handlers/pretty"
	"log/slog"
	"os"
)

const (
	EnvProd = "prod"
	EnvDev  = "dev"
)

func SetupLogger(env string, debug bool) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvProd:

		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log = PrettyLogHandler(env, debug)
	default:

		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func PrettyLogHandler(env string, debug bool) *slog.Logger {
	opts := pretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	logger := slog.New(opts.NewPrettyHandler(os.Stdout))
	if debug == true {
		// logger = logger.With(slog.String("env", env))
	}
	return logger
}
