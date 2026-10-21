package main

import (
	"context"
	"golang_restapi/internal/logger"
	app "golang_restapi/internal/transaction"
	"golang_restapi/internal/transaction/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	l := logger.New(&cfg.Base)
	logger.DumpConfig(l, cfg.Redacted())
	application := app.New(&l, cfg)
	if err := application.Run(ctx); err != nil {
		l.Fatal().Err(err).Msg("application.Run failed")
	}
}
