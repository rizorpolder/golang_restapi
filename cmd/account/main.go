package main

import (
	"context"
	app "golang_restapi/internal/account"
	"golang_restapi/internal/account/config"
	"golang_restapi/internal/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	//_ "golang_restapi/docs"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:50052
// @BasePath /
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

// PingExample godoc
// @Summary Проверка доступности сервиса
// @Description Возвращает pong
// @Tags health
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingExample(c *gin.Context) {
	c.String(200, "pong")
}
