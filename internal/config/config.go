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
