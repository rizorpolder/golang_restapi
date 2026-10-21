package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	Host        string `env:"HTTP_HOST,required" json:"http_host" envDefault:"localhost"`
	Port        int    `env:"HTTP_PORT,required" json:"http_port" envDefault:"9000"`
	GrpcHost    string `env:"GRPC_HOST,required" json:"grpc_host" envDefault:"localhost"`
	GrpcPort    int    `env:"GRPC_PORT,required" json:"grpc_port" envDefault:"50050"`
	JWTSecret   string `env:"JWT_SECRET,required,notEmpty" json:"jwt_secret"`
	AccountAddr string `env:"ACCOUNT_ADDR,required" json:"account_addr"`
	AuthAddr    string `env:"AUTH_ADDR,required" json:"auth_addr"`
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/gateway.env"
	}
	return config.Load[Config](path)
}

func (c Config) Redacted() Config {
	c.JWTSecret = "***"

	return c
}
