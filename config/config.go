package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Environment       string   `env:"ENVIRONMENT" envDefault:"dev"`
	Port              string   `env:"PORT" envDefault:"8081"`
	AllowedOrigins    []string `env:"ALLOW_ORIGINS" envSeparator:"," envDefault:"*"`
	GoogleCloudConfig GoogleCloudConfig
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
