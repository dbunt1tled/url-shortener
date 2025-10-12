package main

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/dbunt1tled/url-shortener/cmd/urlshortener/transport"
	"github.com/dbunt1tled/url-shortener/internal/app/container"
	cfg "github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/lib/http-server/routes"
	"github.com/dbunt1tled/url-shortener/internal/lib/validator"
	"github.com/dbunt1tled/url-shortener/storage/mysql"
	_ "go.uber.org/automaxprocs"
)

func main() {
	env := cfg.LoadConfig()
	logger := cfg.LoadLogger(env.Env, env.LogLevel)
	// localeBundle := cfg.SetupLocale()
	storage := mysql.GetInstance()
	srvOpts := []config.Option{
		server.WithHostPorts(env.HTTPAddress),
		server.WithKeepAlive(true),
		server.WithReadTimeout(env.HTTPTimeout),
		server.WithWriteTimeout(env.HTTPTimeout),
		server.WithIdleTimeout(env.HTTPIdleTimeout),
		server.WithCustomValidator(validator.New()),
		server.WithExitWaitTime(3 * time.Second), //nolint:mnd // 3sec enough
	}
	srv := transport.NewServer(srvOpts...)

	srv.OnShutdown = append(srv.OnShutdown, func(ctx context.Context) {
		logger.Info("㋡ Quit: closing database connection")
		err := storage.Close()
		if err != nil {
			logger.Error("Error closing storage", err)
		}
		<-ctx.Done()
		logger.Warn("Exit timeout!")
	})
	di, err := container.Build()
	if err != nil {
		panic(err)
	}
	err = di.Provide(func() *server.Hertz {
		return srv
	})
	if err != nil {
		panic(err)
	}

	err = di.Provide(routes.NewRouter)
	if err != nil {
		panic(err)
	}

	err = di.Invoke(func(r *routes.Router) {
		r.Register()
	})
	if err != nil {
		panic(err)
	}

	logger.Debug("Start listening on address: " + env.HTTPAddress)
	srv.Spin()
}
