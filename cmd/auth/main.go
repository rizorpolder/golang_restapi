package auth

import (
	"context"
	app "golang_restapi/internal/account"
	"golang_restapi/internal/config"
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

	l := logger.New(cfg)

	application := app.New(&l, cfg)
	if err := application.Run(ctx); err != nil {
		l.Fatal().Err(err).Msg("Application run failed")
	}
}
