package config

import "github.com/caarlos0/env/v9"

type GoogleCloudConfig struct {
	ProjectID    string `env:"PROJECT_ID"`
	PubsubConfig PubSubConfig
}

func NewGoogleCloudConfig() (*GoogleCloudConfig, error) {
	cfg := &GoogleCloudConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
