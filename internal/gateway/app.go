package app

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"net"

	_ "github.com/lib/pq"
	accountpb "golang_restapi/contracts/account/go"
	authbp "golang_restapi/contracts/auth/go"
	gatewaypb "golang_restapi/contracts/gateway/go"
	"golang_restapi/internal/gateway/account"
	"golang_restapi/internal/gateway/auth"
	"golang_restapi/internal/gateway/config"
	"golang_restapi/internal/gateway/server"
	"golang_restapi/internal/gateway/service"
	"golang_restapi/internal/interceptor"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg    *config.Config
	logger *zerolog.Logger

	service        *service.GatewayService
	authService    *auth.Service
	accountService *account.Service

	server     *server.Server
	grpcServer *grpc.Server
}

func NewApp(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	gatewayServer, err := a.getGatewayServer()
	if err != nil {
		return fmt.Errorf("failed get gateway server: %w", err)
	}

	a.grpcServer = a.getGrpcServer(gatewayServer)

	listenAddr := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		a.logger.Fatal().Err(err).Msg("failed to listen")
		return err
	}
	a.logger.Info().Msg("gateway server listening")

	serveErrCh := make(chan error)
	go func() {
		serveErrCh <- a.grpcServer.Serve(listener)
	}()
	select {
	case <-ctx.Done():
		{
			a.grpcServer.GracefulStop()
			return ctx.Err()
		}
	case err := <-serveErrCh:
		if err != nil {
			a.logger.Fatal().Err(err).Msg("failed to serve")
		}
		return err
	}
}

func (a *App) getGatewayServer() (*server.Server, error) {
	if a.server == nil {
		service, err := a.getGatewayService()
		if err != nil {
			return nil, fmt.Errorf("failed get server %w", err)
		}
		a.server = server.New(*service, a.logger)
	}
	return a.server, nil
}

func (a *App) getGatewayService() (*service.GatewayService, error) {
	if a.service == nil {
		authService, err := a.getAuthService()

		if err != nil {
			return nil, fmt.Errorf("failed to get auth service: %w", err)
		}

		accountService, err := a.getAccountService()
		if err != nil {
			return nil, fmt.Errorf("failed to get account service: %w", err)
		}
		a.service = service.New(accountService, authService, a.logger)
	}
	return a.service, nil
}

func (a *App) getAuthService() (*auth.Service, error) {
	if a.authService == nil {
		conn, err := grpc.NewClient(a.cfg.AuthAddr, grpc.WithTransportCredentials(
			insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed create auth service: %w", err)
		}
		client := authbp.NewAuthClient(conn)
		a.authService = auth.NewService(client)
	}
	return a.authService, nil
}

func (a *App) getAccountService() (*account.Service, error) {
	if a.accountService == nil {
		conn, err := grpc.NewClient(
			a.cfg.AccountAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc client: %w", err)
		}
		client := accountpb.NewAccountClient(conn)
		a.accountService = account.NewService(client)
	}
	return a.accountService, nil
}

func (a *App) getGrpcServer(srv *server.Server) *grpc.Server {
	jwtInterceptor := interceptor.NewJWTInterceptor(a.cfg.JWTSecret)
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(jwtInterceptor.UnaryInterceptor()))
	gatewaypb.RegisterGatewayServer(grpcSrv, srv)
	return grpcSrv
}

func (a *App) Close() error {
	if a.grpcServer == nil {
		a.grpcServer.GracefulStop()
	}
	return nil
}
