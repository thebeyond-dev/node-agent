package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MigrationsPath string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	DatabasePath   string `env:"DATABASE_PATH" env-default:"./data/ips.db"`
	Port           string `env:"PORT" env-default:"4080"`
	LogLevel       string `env:"LOG_LEVEL" env-default:"info"`
	NodeIP         string `env:"NODE_IP" env-required:"true"`
	NetworkCIDR    string `env:"NETWORK_CIDR" env-default:"10.10.0.0/20"`
	ClientDNS      string `env:"CLIENT_DNS" env-default:"1.1.1.1,8.8.8.8"`
	AuthSecret     string `env:"AUTH_SECRET" env-required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &cfg, nil
}
