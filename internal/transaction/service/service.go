package service

import (
	"context"
	"github.com/rs/zerolog"
	"golang_restapi/internal/transaction/model"
)

type TransactionService struct {
	repo   Repository
	logger *zerolog.Logger
}

func New(repo Repository, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{
		repo:   repo,
		logger: logger,
	}
}

type Repository interface {
	GetTransactions(context.Context, model.GetTransactionParams) ([]model.Transaction, error)
	GetTransactionDetails(context.Context, uint64) (model.TransactionDetails, error)
	Withdraw(context.Context, model.WithdrawParams) (model.TransactionDetails, error)

	Deposit(context.Context, model.DepositParams) (model.TransactionDetails, error)
	Transfer(context.Context, model.TransferParams) (model.TransactionDetails, error)
}

func (s *TransactionService) GetTransactionsWithDetails(ctx context.Context, params model.GetTransactionParams) ([]model.TransactionDetails, error) {
	transactions, err := s.repo.GetTransactions(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get transactions")
		return nil, err
	}

	var transactionDetails []model.TransactionDetails
	for _, transaction := range transactions {
		details, err := s.repo.GetTransactionDetails(ctx, transaction.ID)
		if err != nil {
			s.logger.Error().Err(err).Uint64("transaction_id", transaction.ID).Msg("failed to get transaction details")
			continue
		}
		transactionDetails = append(transactionDetails, details)
	}
	return transactionDetails, nil
}

func (s *TransactionService) Deposit(ctx context.Context, userId uint64, amount float64) (model.TransactionDetails, error) {
	params := model.DepositParams{
		UserID: userId,
		Amount: amount,
	}
	return s.repo.Deposit(ctx, params)
}

func (s *TransactionService) Withdraw(ctx context.Context, accountId uint64, amount float64) (model.TransactionDetails, error) {
	params := model.WithdrawParams{
		AccountID: accountId,
		Amount:    amount,
	}
	return s.repo.Withdraw(ctx, params)
}

func (s *TransactionService) Transfer(ctx context.Context, userId, recipient uint64, amount float64) (model.TransactionDetails, error) {
	params := model.TransferParams{
		UserID:    userId,
		Recipient: recipient,
		Amount:    amount,
	}
	return s.repo.Transfer(ctx, params)
}
