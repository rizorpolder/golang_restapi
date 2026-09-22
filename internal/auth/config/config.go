package config

import "golang_restapi/internal/config"

type Config struct {
	config.Base
	DBDSN                 string `env:"DB_DSN,required"`
	JWTSecret             string `env:"JWT_SECRET,required"`
	AccessTokenTTLMinutes int    `env:"ACCESS_TOKEN_TTL_MINUTES" envDefault:"60"`
	RefreshTokenTTLDays   int    `env:"REFRESH_TOKEN_TTL_DAYS" envDefault:"30"`
}
