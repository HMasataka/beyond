package config

import "github.com/caarlos0/env/v9"

type Firebase struct {
	Credentials string `env:"FIREBASE_CREDENTIALS"`
}

func NewFirebaseConfig() (*Firebase, error) {
	cfg := &Firebase{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
