package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" required:"true"`
	AppEnv      string `env:"APP_ENV" required:"true" envDefault:"development"`
	Host        string `env:"HTTP_HOST" required:"true" envDefault:"localhost"`
	Port        string `env:"HTTP_PORT" required:"true" envDefault:"9000"`
	Loglevel    string `env:"LOG_LEVEL" required:"true" envDefault:"info"`

	DbDns string `env:"DB_DNS" required:"true"`

	JwtSecret             string `env:"JWT_SECRET" json:"jwt_secret" required:"true"`
	AccessTokenTTLMinutes int    `env:"ACCESS_TOKEN_TTL_MINUTES" json:"access_ttl_min" required:"true" default:"60"`
	RefreshTokenTTLDays   int    `env:"REFRESH_TOKEN_TTL_DAYS" json:"refresh_ttl_days" required:"true" default:"30"`
}

func Load() (*Config, error) {

	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: no .env file found")
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
