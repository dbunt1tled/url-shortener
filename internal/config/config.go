package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	AppName     string `env:"APP_NAME" env-required:"true"`
	Env         string `env:"ENV" env-default:"dev"`
	Debug       bool   `env:"DEBUG" env-default:"false"`
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	HTTPServer  HTTPServer
}
type HTTPServer struct {
	Address     string `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080" env-required:"true"`
	Timeout     string `env:"HTTP_SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout string `env:"HTTP_SERVER_IDLE_TIMEOUT" env-required:"true"`
}

func MustLoadConfig() *Config {
	var cfg Config
	configPath := ".env"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("Error load config enviroment: %s", err)
		}
	} else {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Error load config file enviroment: %s", err)
		}
	}

	return &cfg
}
