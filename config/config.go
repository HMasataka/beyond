package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Environment    string   `env:"ENVIRONMENT" envDefault:"dev"`
	Port           string   `env:"PORT" envDefault:"8080"`
	AllowedOrigins []string `env:"ALLOW_ORIGINS" envSeparator:"," envDefault:"*"`
	GoogleCloud    GoogleCloud
	Firebase       Firebase
	MySQL          MySQL
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
