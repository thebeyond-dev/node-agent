package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MigrationsPath string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	DatabasePath   string `env:"DATABASE_PATH" env-default:"./data/ips.db"`
	Port           string `env:"PORT" env-default:"2000"`
	LogLevel       string `env:"LOG_LEVEL" env-default:"info"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &cfg, nil
}
