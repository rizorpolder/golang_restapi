package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	Host        string `env:"HTTP_HOST" json:"http_host" required:"true" default:"localhost"`
	Port        int    `env:"HTTP_PORT" json:"http_port" required:"true" default:"9000"`
	GrpcHost    string `env:"GRPC_HOST" json:"grpc_host" required:"true" default:"localhost"`
	GrpcPort    int    `env:"GRPC_PORT" json:"grpc_port" required:"true" default:"50050"`
	JWTSecret   string `env:"JWT_SECRET" json:"jwt_secret" required:"true"`
	AccountAddr string `env:"ACCOUNT_ADDR" json:"account_addr" required:"true"` // "account:50051"
	AuthAddr    string `env:"AUTH_ADDR" json:"auth_addr" required:"true"`       // "auth:50052"

}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/gateway.env"
	}
	return config.Load[Config](path)
}
