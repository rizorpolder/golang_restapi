package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Base struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load[T any](envFile string) (*T, error) {
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	cfg := new(T)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
