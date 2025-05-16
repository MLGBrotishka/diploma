package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config - структура для хранения конфигурации
	Config struct {
		App      App      `yaml:"app"`
		Log      Log      `yaml:"log"`
		Postgres Postgres `yaml:"postgres"`
		GRPC     GRPC     `yaml:"grpc"`
		HTTP     HTTP     `yaml:"http"`
		JWT      JWT      `yaml:"jwt"`
	}

	App struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	Postgres struct {
		URL string `yaml:"url" env:"POSTGRES_URL"`
	}

	GRPC struct {
		Port int `yaml:"port" env:"GRPC_PORT"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	JWT struct {
		Secret string        `yaml:"secret" env:"JWT_SECRET"`
		TTL    time.Duration `yaml:"ttl" env:"JWT_TTL"`
	}
)

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
