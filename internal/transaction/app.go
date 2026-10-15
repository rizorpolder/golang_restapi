package transaction

import (
	"github.com/rs/zerolog"
	"golang_restapi/internal/transaction/config"
	"golang_restapi/internal/transaction/repository"
	"golang_restapi/internal/transaction/service"
)

type App struct {
	cfg                *config.Config
	logger             *zerolog.Logger
	transactionRepo    *repository.Repository
	transactionService *service.TransactionService
	transactionServer  *server.Server
}
