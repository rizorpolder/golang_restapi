package config

import "golang_restapi/internal/config"

type Config struct {
	config.Base
	DBDSN string `env:"DB_DSN,required"`
}
