package app

import (
	"github.com/rs/zerolog"
	"golang_restapi/internal/account/server"
	account "golang_restapi/internal/account/service"
	auth "golang_restapi/internal/auth/service"
	"golang_restapi/internal/gateway/config"
	gateway "golang_restapi/internal/gateway/service"
	"google.golang.org/grpc"
)

type App struct {
	cfg            *config.Config
	logger         *zerolog.Logger
	service        *gateway.GatewayService
	authService    *auth.AuthService
	accountService *account.AccountService
	server         *server.Server
	grpcServer     *grpc.Server
}

func NewApp(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

///todo обвязал gateway сервис на интерфейсы, нужн проверить расхождения в методах (в самих сервисах)
///дописать app.go
