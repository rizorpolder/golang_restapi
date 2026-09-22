package app

//
//import (
//	"github.com/rs/zerolog"
//	account "golang_restapi/contracts/account/go"
//	auth "golang_restapi/contracts/auth/go"
//	gateway "golang_restapi/contracts/gateway/go"
//	"golang_restapi/internal/account/server"
//	"golang_restapi/internal/config"
//	"google.golang.org/grpc"
//)
//
//type App struct {
//	cfg            *config.Config
//	logger         *zerolog.Logger
//	service        *gateway.GatewayService
//	authService    *auth.Service
//	accountService *account.Service
//	server         *server.Server
//	grpcServer     *grpc.Server
//}
//
//func NewApp(logger *zerolog.Logger, cfg *config.Config) *App {
//	return &App{
//		cfg:    cfg,
//		logger: logger,
//	}
//}
