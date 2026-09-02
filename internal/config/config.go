package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" required:"true"`
	AppEnv      string `env:"APP_ENV" required:"true" default:"development"`
	Host        string `env:"HOST" required:"true" default:"localhost"`
	Port        string `env:"PORT" required:"true" default:"9000"`
	Loglevel    string `env:"LOG_LEVEL" required:"true" default:"info"`

	DbDns string `env:"DB_DNS" required:"true"`
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
