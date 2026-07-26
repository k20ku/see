package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env        string `env:"SEE_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"SEE_DB_HOST" envDefault:"127.0.0.1"`
	DBPort     string `env:"SEE_DB_PORT" envDefault:"5432"`
	DBUser     string `env:"SEE_DB_USER" envDefault:"see"`
	DBPassword string `env:"SEE_DB_PASSWORD" envDefault:"seedbpass"`
	DBName     string `env:"SEE_DB_NAME" envDefault:"see"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}
