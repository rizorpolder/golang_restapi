package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	DBDSN    string `env:"DB_DSN,required"`
	GRPCHost string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/account.env"
	}
	return config.Load[Config](path)
}

func (c Config) Redacted() Config {
	c.DBDSN = "***"

	return c
}
