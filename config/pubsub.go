package config

import "github.com/caarlos0/env/v9"

type PubSubConfig struct {
	TopicID string `env:"TOPIC_ID"`
}

func NewPubsubConfig() (*PubSubConfig, error) {
	cfg := &PubSubConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
