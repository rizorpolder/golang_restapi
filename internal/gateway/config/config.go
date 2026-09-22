package config

import "golang_restapi/internal/config"

type Config struct {
	config.Base
	JWTSecret   string `env:"JWT_SECRET,required"`
	AccountAddr string `env:"ACCOUNT_ADDR,required"` // "account:50051"
	AuthAddr    string `env:"AUTH_ADDR,required"`    // "auth:50052"
}
