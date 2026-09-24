package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	authpb "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/auth/config"
	"golang_restapi/internal/auth/repository"
	"golang_restapi/internal/auth/server"
	"golang_restapi/internal/auth/service"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
)

type App struct {
	cfg         *config.Config
	l           *zerolog.Logger
	authRepo    *repository.Repository
	authService *service.AuthService
	authServer  *server.Server
	grpcServer  *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		l:   logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	authServer, err := a.getAuthServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth server: %w", err)
	}
	a.grpcServer = getGRPCServer(authServer)
	listenAddr := fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		a.l.Fatal().Err(err).Msgf("failed to listen on %s:%v", listenAddr, err)
		return err
	}
	a.l.Info().Msgf("listening on %s:%v", listenAddr, a.cfg.Port)
	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- a.grpcServer.Serve(lis)
	}()
	select {
	case <-ctx.Done():
		a.Close(ctx)
		return ctx.Err()

	case err := <-serveErrCh:
		if err != nil {
			a.l.Error().Err(err).Msgf("failed to serve: %v", err)
		}
		return err

	}
}

func (a *App) getAuthServer(ctx context.Context) (*server.Server, error) {
	if a.authServer == nil {
		service, err := a.getAuthService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get auth server: %w", err)
		}
		a.authServer = server.New(*service, a.l)
	}
	return a.authServer, nil
}

func (a *App) getAuthService(ctx context.Context) (*service.AuthService, error) {
	if a.authService == nil {
		repo, err := a.getRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}
		a.authService = service.New(*repo, a.cfg, a.l)

	}
	return a.authService, nil
}

func (a *App) getRepository(ctx context.Context) (*repository.Repository, error) {
	if a.authRepo == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	db, err := gorm.Open(postgres.Open(a.cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	a.authRepo = repository.NewRepository(db, a.l)
	return a.authRepo, nil
}

func (a *App) runMigrations(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set postgres dialect: %w", err)
	}
	dbGoose, err := sql.Open("postgres", a.cfg.DBDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := goose.UpContext(ctx, dbGoose, "internal/auth/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func getGRPCServer(authServer *server.Server) *grpc.Server {
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServer(grpcServer, authServer)
	return grpcServer
}

func (a *App) Close(ctx context.Context) error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	return nil
}
