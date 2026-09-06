package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	accountpb "golang_restapi/contracts/account/go"
	"golang_restapi/internal/account/repository"
	"golang_restapi/internal/account/server"
	"golang_restapi/internal/account/service"
	"golang_restapi/internal/config"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"

	_ "github.com/lib/pq"
	_ "golang_restapi/internal/account/migrations"
)

type App struct {
	cfg            *config.Config
	l              *zerolog.Logger
	accountRepo    *repository.Repository
	accountService *service.AccountService
	accountServer  *server.Server
	grpcServer     *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		l:   logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	accountServer, err := a.getAccountServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account server: %w", err)
	}
	a.grpcServer = getGRPCServer(accountServer)

	listenAddr := fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port)

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		a.l.Fatal().Err(err).Msgf("Failed to listen on %s: %v", listenAddr, err)
		return err
	}

	a.l.Info().Msgf("gRPC servcer listening on %v", listenAddr)

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
			a.l.Error().Err(err).Msgf("Fail to serve: %v", err)
		}
		return err
	}
}

func (a *App) getRepository(ctx context.Context) (*repository.Repository, error) {
	if a.accountRepo == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	db, err := gorm.Open(postgres.Open(a.cfg.DbDns), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	a.accountRepo = repository.NewRepository(db, a.l)
	return a.accountRepo, nil
}

func (a *App) runMigrations(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set postgres dialect: %w", err)
	}
	dbGoose, err := sql.Open("postgres", a.cfg.DbDns)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := goose.UpContext(ctx, dbGoose, "internal/account/migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
func (a *App) getAccountService(ctx context.Context) (*service.AccountService, error) {
	if a.accountService == nil {
		repo, err := a.getRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}
		a.accountService = service.New(*repo, a.l)
	}
	return a.accountService, nil
}

func (a *App) getAccountServer(ctx context.Context) (*server.Server, error) {
	if a.accountServer == nil {
		service, err := a.getAccountService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get account server: %w", err)
		}
		a.accountServer = server.New(*service, a.l)
	}
	return a.accountServer, nil
}

func getGRPCServer(accountServer *server.Server) *grpc.Server {
	grpcServer := grpc.NewServer()
	accountpb.RegisterAccountServer(grpcServer, accountServer)
	return grpcServer
}

func (a *App) Close(ctx context.Context) error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	return nil
}
