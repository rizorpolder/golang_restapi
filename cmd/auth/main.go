package main

import (
	"context"
	app "golang_restapi/internal/auth"
	"golang_restapi/internal/auth/config"
	"golang_restapi/internal/logger"
)

// @title Auth Service
// @version 1.0
// @description Auth Service
// @host localhost:50051
// @BasePath /
func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(&cfg.Base)
	logger.DumpConfig(log, cfg.Redacted())

	application := app.New(&log, cfg)
	if err := application.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Application run failed")
	}
}
