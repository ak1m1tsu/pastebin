package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Postgres `yaml:"postgres"`
		Redis    `yaml:"redis"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"PORT"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	Postgres struct {
		DSN string `yaml:"dsn" env:"PG_DSN"`
	}

	Redis struct {
		DSN string `yaml:"dsn" env:"REDIS_DSN"`
	}
)

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
