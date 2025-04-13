package config

import "github.com/caarlos0/env/v9"

type MySQL struct {
	Net      string `env:"MYSQL_NET" envDefault:"tcp"`
	DB       string `env:"MYSQL_DB" envDefault:"db"`
	Password string `env:"MYSQL_PASSWORD" envDefault:"password"`
	User     string `env:"MYSQL_USER" envDefault:"user"`
	Host     string `env:"MYSQL_HOST" envDefault:"localhost"`
	Port     int    `env:"MYSQL_PORT" envDefault:"3306"`
}

func NewMySQLConfig() (*MySQL, error) {
	cfg := &MySQL{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
