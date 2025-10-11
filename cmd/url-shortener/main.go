package main

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/netpoll"
	"github.com/dbunt1tled/url-shortener/internal/app/container"
	"github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/routes"
	"github.com/dbunt1tled/url-shortener/internal/lib/validator"
	"github.com/dbunt1tled/url-shortener/storage/mysql"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cfg := config.LoadConfig()
	logger := config.LoadLogger(cfg.Env, cfg.LogLevel)
	// localeBundle := config.SetupLocale()
	storage := mysql.GetInstance()
	srv := server.New(
		server.WithTransport(netpoll.NewTransporter),
		server.WithHostPorts(cfg.HTTPAddress),
		server.WithKeepAlive(true),
		server.WithReadTimeout(cfg.HTTPTimeout),
		server.WithWriteTimeout(cfg.HTTPTimeout),
		server.WithIdleTimeout(cfg.HTTPIdleTimeout),
		server.WithCustomValidator(validator.New()),
		server.WithExitWaitTime(3*time.Second), //nolint:mnd // 3sec enough
	)

	srv.OnShutdown = append(srv.OnShutdown, func(ctx context.Context) {
		logger.Info("㋡ Quit: closing database connection")
		err := storage.Close()
		if err != nil {
			logger.Error("Error closing storage", err)
		}
		<-ctx.Done()
		logger.Warn("Exit timeout!")
	})
	container, err := container.Build()
	if err != nil {
		panic(err)
	}
	err = container.Provide(func() *server.Hertz {
		return srv
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(routes.NewRouter)
	if err != nil {
		panic(err)
	}

	err = container.Invoke(func(r *routes.Router) {
		r.Register()
	})
	if err != nil {
		panic(err)
	}

	logger.Debug("Start listening on address: " + cfg.HTTPAddress)
	srv.Spin()
}
