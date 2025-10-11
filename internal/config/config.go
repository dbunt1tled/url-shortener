package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	HTTPTimeout        = 10 * time.Second
	HTTPIdleTimeout    = 120 * time.Second
	JWTAccessLifeTime  = 3600 * time.Second
	JWTRefreshLifeTime = 7200 * time.Second
)

type Config struct {
	AppName                    string        `mapstructure:"app_name"`
	Env                        string        `mapstructure:"env"`
	LogLevel                   slog.Level    `mapstructure:"log_level"`
	LogLevelStatus             int           `mapstructure:"logger_level_status"`
	DatabaseDSN                string        `mapstructure:"database_dsn"`
	HTTPAddress                string        `mapstructure:"http_server_address"`
	HTTPTimeout                time.Duration `mapstructure:"http_server_timeout"`
	HTTPIdleTimeout            time.Duration `mapstructure:"http_server_idle_timeout"`
	AccessControlAllowHeaders  string        `mapstructure:"access_control_allow_headers"`
	AccessControlExposeHeaders string        `mapstructure:"access_control_expose_headers"`
	AccessControlAllowMethods  string        `mapstructure:"access_control_allow_methods"`
	AccessControlAllowOrigin   string        `mapstructure:"access_control_allow_origin"`
	JWTPublicKey               string        `mapstructure:"jwt_public_key"`
	JWTPrivateKey              string        `mapstructure:"jwt_private_key"`
	JWTAlgorithm               string        `mapstructure:"jwt_token_algorithm"`
	AccessLifeTime             time.Duration `mapstructure:"token_access_life_time_seconds"`
	RefreshLifeTime            time.Duration `mapstructure:"token_refresh_life_time_seconds"`
	BaseURL                    string        `mapstructure:"base_url"`
}

var (
	configInstance *Config   //nolint:gochecknoglobals // singleton
	cm             sync.Once //nolint:gochecknoglobals // singleton
)

func LoadConfig() *Config {
	cm.Do(func() {
		var err error
		configInstance, err = load()
		if err != nil {
			panic("Config is not initialized: " + err.Error())
		}
	})
	return configInstance
}

func load() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("env", "dev")
	v.SetDefault("log_level", 0)
	v.SetDefault("logger_level_status", consts.StatusBadRequest)
	v.SetDefault("http_server_address", ":8080")
	v.SetDefault("http_server_timeout", HTTPTimeout)
	v.SetDefault("http_server_idle_timeout", HTTPIdleTimeout)
	v.SetDefault("access_control_allow_headers", "Origin, Content-Type, Accept, Authorization")
	v.SetDefault("access_control_expose_headers", "Content-Length, Content-Type")
	v.SetDefault("access_control_allow_methods", "GET, POST, PUT, DELETE, OPTIONS")
	v.SetDefault("access_control_allow_origin", "*")
	v.SetDefault("jwt_public_key", "")
	v.SetDefault("jwt_private_key", "")
	v.SetDefault("jwt_token_algorithm", "HS256")
	v.SetDefault("token_access_life_time_seconds", JWTAccessLifeTime)
	v.SetDefault("token_refresh_life_time_seconds", JWTRefreshLifeTime)

	if _, err := os.Stat(".env"); err == nil {
		v.SetConfigFile(".env")
		v.SetConfigType("env")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.AppName == "" {
		return errors.New("App name is required")
	}
	if cfg.DatabaseDSN == "" {
		return errors.New("database dsn is required")
	}
	if cfg.HTTPAddress == "" {
		return errors.New("HTTP Adders:Port is required")
	}
	if cfg.JWTPrivateKey == "" {
		return errors.New("Private Key is required")
	}
	if cfg.JWTAlgorithm == "" || cfg.JWTAlgorithm == "None" {
		return errors.New("Algorithm is required")
	}
	if cfg.BaseURL == "" {
		return errors.New("Base URL is required")
	}
	return nil
}
