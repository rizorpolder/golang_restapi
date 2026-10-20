package transaction

import (
	"context"
	"database/sql"
	"fmt"
	accountpb "golang_restapi/contracts/account/go"
	transactionpb "golang_restapi/contracts/transaction/go"
	kafkaConfig "golang_restapi/internal/kafka/config"
	kafka "golang_restapi/internal/kafka/transaction"
	"golang_restapi/internal/transaction/account"
	"golang_restapi/internal/transaction/config"
	"golang_restapi/internal/transaction/repository"
	"golang_restapi/internal/transaction/server"
	"golang_restapi/internal/transaction/service"
	"net"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

const migrationsDir = "internal/transaction/migrations"

type App struct {
	cfg                *config.Config
	logger             *zerolog.Logger
	transactionRepo    *repository.Repository
	transactionService *service.TransactionService
	transactionServer  *server.Server
	grpcServer         *grpc.Server

	accountService *account.Service
	kafkaClient    *kafka.Kafka
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	transactionServer, err := a.getTransactionServer(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction server: %w", err)
	}
	a.grpcServer = getGrpcServer(transactionServer)
	listenAdd := net.JoinHostPort(a.cfg.GRPCHost, a.cfg.GRPCPort)
	lis, err := net.Listen("tcp", listenAdd)
	if err != nil {
		a.logger.Fatal().Err(err).Msg("could not listen on " + listenAdd)
		return err
	}
	a.logger.Log().Msgf("listening on %s", lis.Addr())
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
			a.logger.Fatal().Err(err).Msg("could not start server")
		}
		return err

	}
}

func (a *App) getTransactionServer(ctx context.Context) (*server.Server, error) {
	if a.transactionService == nil {
		svc, err := a.getTransactionService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction service: %w", err)
		}
		a.transactionServer = server.New(svc, a.logger)
	}
	return a.transactionServer, nil
}

func (a *App) getTransactionService(ctx context.Context) (*service.TransactionService, error) {
	if a.transactionService == nil {
		repo, err := a.getTransactionRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction repository: %w", err)
		}

		accountSvc, err := a.getAccountService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get account service: %w", err)
		}

		kafkaClient, err := a.getKafkaClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get kafka client: %w", err)
		}

		a.transactionService = service.New(repo, accountSvc, kafkaClient, a.logger)
		a.kafkaClient.Subscribe(ctx, "transaction_response", a.transactionService.HandleAccountResponse)
	}
	return a.transactionService, nil
}

func (a *App) getKafkaClient(ctx context.Context) (*kafka.Kafka, error) {
	if a.kafkaClient == nil {
		producerCfg := kafkaConfig.DefaultProducerConfig(a.cfg.KafkaBrokers)
		producer := kafka.NewProducer(producerCfg, a.logger)

		kafkaClient := kafka.New(producer, a.cfg.KafkaBrokers, a.cfg.KafkaGroupID, a.logger)
		a.kafkaClient = kafkaClient
		a.logger.Info().Msg("kafka client created")
	}
	return a.kafkaClient, nil
}

func (a *App) getAccountService(ctx context.Context) (*account.Service, error) {
	if a.accountService == nil {
		conn, err := grpc.NewClient(a.cfg.AccountGrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("could not create grpc client: %w", err)
		}
		client := accountpb.NewAccountClient(conn)
		a.accountService = account.New(client)
	}
	return a.accountService, nil
}

func getGrpcServer(transactionServer *server.Server) *grpc.Server {
	grpcServer := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcServer, transactionServer)
	return grpcServer
}

func (a *App) getTransactionRepository(ctx context.Context) (*repository.Repository, error) {
	if a.transactionRepo == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}
	db, err := gorm.Open(postgres.Open(a.cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	a.transactionRepo = repository.NewRepository(db, a.logger)
	return a.transactionRepo, nil
}

func (a *App) runMigrations(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	dbGoose, err := sql.Open("postgres", a.cfg.DBDSN)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	if err := goose.UpContext(ctx, dbGoose, migrationsDir); err != nil {
		fmt.Printf("failed to run migrations: %s\n", err)
	}
	return nil

}
func (a *App) Close(ctx context.Context) error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	return nil
}
