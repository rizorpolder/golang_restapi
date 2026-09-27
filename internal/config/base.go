package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"io/fs"
)

type Base struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load[T any](envFile string) (*T, error) {
	for _, f := range []string{envFile, ".env"} {
		if err := godotenv.Load(f); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("load %s: %w", f, err)
		}
	}

	cfg := new(T)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
