// Package config отвечает за чтение и разбор конфигурации приложения.
package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config is a struct for storing configuration.
	Config struct {
		// App contains application settings.
		App App `yaml:"app"`
		// Log contains logging settings.
		Log Log `yaml:"log"`
		// Postgres contains Postgres database settings.
		Postgres Postgres `yaml:"postgres"`
		// GRPC contains gRPC server settings.
		GRPC GRPC `yaml:"grpc"`
		// HTTP contains HTTP server settings.
		HTTP HTTP `yaml:"http"`
		// Auth contains auth service settings.
		Auth Auth `yaml:"auth"`
	}

	// App contains application settings.
	App struct {
		// Name is the application name.
		Name string `yaml:"name" env:"APP_NAME"`
		// Version is the application version.
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// Log contains logging settings.
	Log struct {
		// Level is the logging level.
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	// Postgres contains Postgres database settings.
	Postgres struct {
		// URL is the Postgres connection URL.
		URL string `yaml:"url" env:"POSTGRES_URL"`
	}

	// GRPC contains gRPC server settings.
	GRPC struct {
		// Port is the gRPC server port.
		Port int `yaml:"port" env:"GRPC_PORT"`
	}

	// HTTP contains HTTP server settings.
	HTTP struct {
		// Port is the HTTP server port.
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	// Auth contains auth service settings.
	Auth struct {
		// GRPC contains gRPC client settings.
		GRPC AuthGRPC `yaml:"grpc"`
	}

	// GRPC contains gRPC client settings.
	AuthGRPC struct {
		// Addr is the gRPC server address.
		Addr string `yaml:"addr" env:"AUTH_GRPC_ADDR"`
	}
)

// NewConfig creates a new Config instance.
func NewConfig(env string) (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(env, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read env: %w", err)
	}
	return cfg, nil
}
