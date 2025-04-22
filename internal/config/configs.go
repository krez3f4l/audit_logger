package config

import envcfg "github.com/kelseyhightower/envconfig"

type Config struct {
	DB     Mongo
	Server Server
}

type Mongo struct {
	URI      string
	Username string
	Password string
	Database string
}

type Server struct {
	Port int
}

func New() (*Config, error) {
	cfg := new(Config)

	if err := envcfg.Process("db", &cfg.DB); err != nil {
		return nil, err
	}

	if err := envcfg.Process("server", &cfg.Server); err != nil {
		return nil, err
	}

	return cfg, nil
}
