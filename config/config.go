package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config - структура для хранения конфигурации
	Config struct {
		App      int      `yaml:"app"`
		Log      Log      `yaml:"log"`
		Postgres Postgres `yaml:"postgres"`
		GRPC     GRPC     `yaml:"grpc"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	Postgres struct {
		URL string `yaml:"url"`
	}

	GRPC struct {
		Port int `yaml:"port"`
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
