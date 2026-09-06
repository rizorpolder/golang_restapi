package main

import (
	"context"
	"github.com/gin-gonic/gin"
	app "golang_restapi/internal/account"
	"golang_restapi/internal/config"
	"golang_restapi/internal/logger"
	//_ "golang_restapi/docs"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:8080
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
