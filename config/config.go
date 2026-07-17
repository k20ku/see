package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env  string `env:"SEE_ENV" envDefault:"dev"`
	Port int    `env:"PORT" envDefault:"80"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %v", err)
	}
	return cfg, nil
}
