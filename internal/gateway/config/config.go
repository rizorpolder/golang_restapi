package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	JWTSecret   string `env:"JWT_SECRET,required"`
	AccountAddr string `env:"ACCOUNT_ADDR,required"` // "account:50051"
	AuthAddr    string `env:"AUTH_ADDR,required"`    // "auth:50052"
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "/configs/gateway.env"
	}
	return config.Load[Config](path)
}
