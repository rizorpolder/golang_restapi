package main

import (
	"golang_restapi/internal/config"
	"golang_restapi/internal/logger"
	"log"

	_ "golang_restapi/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	l := logger.New(cfg)

	//db, err := gorm.Open(postgres.Open(cfg.DbDns), &gorm.Config{})
	//if err != nil {
	//	l.Error().Msgf("Error connecting to DB: %v", err)
	//	return
	//}
	//
	//l.Info().Msgf("Connected to DB")
	//
	//repo := repository.NewRepository(db, &l)
	//_ = repo

	router := gin.Default()
	err = router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		l.Error().Msgf("Error setting trusted proxies: %v", err)
		return
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", PingExample)
	err = router.Run(":8080")

	if err != nil {
		return
	}

	l.Info().Msgf("Service starting up")
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
