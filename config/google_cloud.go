package config

import "github.com/caarlos0/env/v9"

type GoogleCloud struct {
	ProjectID    string `env:"PROJECT_ID"`
	PubsubConfig PubSubConfig
}

func NewGoogleCloudConfig() (*GoogleCloud, error) {
	cfg := &GoogleCloud{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
