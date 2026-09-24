package config

import (
	"golang_restapi/internal/config"
	"os"
)

type Config struct {
	config.Base
	DBDSN                 string `env:"DB_DSN,required"`
	JWTSecret             string `env:"JWT_SECRET,required"`
	AccessTokenTTLMinutes int    `env:"ACCESS_TOKEN_TTL_MINUTES" envDefault:"60"`
	RefreshTokenTTLDays   int    `env:"REFRESH_TOKEN_TTL_DAYS" envDefault:"30"`
}

func Load() (*Config, error) {
	path := os.Getenv("ENV_FILE")
	if path == "" {
		path = "configs/auth.env"
	}
	return config.Load[Config](path)
}

func (c Config) Redacted() Config {
	c.DBDSN = "***"
	c.JWTSecret = "***"
	c.DBDSN = "***"
	return c
}
