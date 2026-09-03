package main

import (
	"fmt"
	"golang_restapi/internal/config"
	"golang_restapi/internal/logger"
	"golang_restapi/internal/repository/account"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	l := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DbDns), &gorm.Config{})
	if err != nil {
		l.Error().Msgf("Error connecting to DB: %v", err)
		return
	}

	l.Info().Msgf("Connected to DB")

	repo := account.NewRepository(db, &l)
	_ = repo

	//router := gin.Default()

	listenAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		l.Error().Msgf("Fail to listen %s: %v", listenAddr, err)
		return
	}

	grpcServer := grpc.NewServer()

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)
	reflection.Register(grpcServer)

	l.Info().Msgf("gRPC servcer listening on %v", listenAddr)

	if err := grpcServer.Serve(lis); err != nil {
		l.Error().Msgf("Fail to serve: %v", err)
		return
	}
	//err = router.SetTrustedProxies([]string{"127.0.0.1"})
	//if err != nil {
	//	l.Error().Msgf("Error setting trusted proxies: %v", err)
	//	return
	//}
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.GET("/ping", PingExample)
	//err = router.Run(":8080")
	//
	//if err != nil {
	//	return
	//}

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
