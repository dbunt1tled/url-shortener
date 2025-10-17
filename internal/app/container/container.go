package container

import (
	"github.com/dbunt1tled/url-shortener/internal/app/shorturl"
	"github.com/dbunt1tled/url-shortener/internal/config"
	"github.com/dbunt1tled/url-shortener/internal/lib/hasher"
	"github.com/dbunt1tled/url-shortener/internal/lib/locale"
	"github.com/dbunt1tled/url-shortener/storage/mysql"
	"go.uber.org/dig"
)

func Build() (*dig.Container, error) {
	c := dig.New()
	err := c.Provide(config.LoadConfig)
	if err != nil {
		return nil, err
	}

	err = c.Provide(func(cfg *config.Config) *config.AppLogger {
		return config.LoadLogger(cfg.Env, cfg.LogLevel)
	})

	if err != nil {
		return nil, err
	}

	err = c.Provide(locale.SetupLocale)
	if err != nil {
		return nil, err
	}

	err = c.Provide(func(cfg *config.Config) *hasher.Hasher {
		h, err := hasher.NewHasher(cfg.JWTAlgorithm, cfg.JWTPublicKey, cfg.JWTPrivateKey)
		if err != nil {
			panic(err)
		}
		return h
	})

	if err != nil {
		return nil, err
	}

	err = c.Provide(mysql.GetInstance)
	if err != nil {
		return nil, err
	}

	err = c.Provide(shorturl.NewURLRepository)
	if err != nil {
		return nil, err
	}

	err = c.Provide(func(
		urlRepository *shorturl.URLRepository,
		hasher *hasher.Hasher,
		cfg *config.Config,
	) *shorturl.URLService {
		return shorturl.NewURLService(urlRepository, hasher, cfg.BaseURL)
	})
	if err != nil {
		return nil, err
	}

	err = c.Provide(shorturl.NewURLHandler)
	if err != nil {
		return nil, err
	}

	return c, nil
}
