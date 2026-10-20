package service

import (
	"context"
	"github.com/rs/zerolog"
	"golang_restapi/internal/transaction/account"
	"golang_restapi/internal/transaction/model"
)

type TransactionService struct {
	repo   Repository
	logger *zerolog.Logger

	accountService account.Service
}

func New(repo Repository, accountService account.Service, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{
		repo:           repo,
		accountService: accountService,
		logger:         logger,
	}
}

type Repository interface {
	GetTransactions(context.Context, model.GetTransactionParams) ([]model.Transaction, error)
	GetTransactionDetails(context.Context, uint64) (model.TransactionDetails, error)
	Withdraw(context.Context, model.WithdrawParams) (model.TransactionDetails, error)
	Deposit(context.Context, model.DepositParams) (model.TransactionDetails, error)
	Transfer(context.Context, model.TransferParams) (model.TransactionDetails, error)
	UpdateTransactionStatus(context.Context, uint64, model.TransactionStatus) error
}

type AccountService interface {
	GetBalance(ctx context.Context, userId uint64) (int64, error)
	Deposit(ctx context.Context, userId uint64, amount int64, operationId uint64) (int64, error)
	Withdraw(ctx context.Context, userID uint64, amount int64, operationId uint64) (int64, error)
	Transfer(ctx context.Context, userID, recipient uint64, amount int64, operationId uint64) (string, error)
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
	res, err := s.repo.Deposit(ctx, params) //создаем операцию
	if err != nil {
		s.logger.Error().Err(err).Uint64("user_id", userId).Msg("failed to create deposit transaction")
		return model.TransactionDetails{}, err
	}

	//создаем транзакцию
	_, err = s.accountService.Deposit(ctx, userId, int64(amount*100), res.Transaction.ID)
	if err != nil {
		s.logger.Error().Err(err).Uint64("user_id", userId).Uint64("transaction_id", res.Transaction.ID).Msg("failed to deposit")
		updateErr := s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusFailed)
		if updateErr != nil {
			s.logger.Error().Err(updateErr).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update transaction status")
		}
		return model.TransactionDetails{}, updateErr
	}
	err = s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusCompleted)
	if err != nil {
		s.logger.Error().Err(err).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update transaction status")
		return model.TransactionDetails{}, err
	}
	res.Transaction.Status = model.TransactionStatusCompleted
	return res, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountId uint64, amount float64) (model.TransactionDetails, error) {
	params := model.WithdrawParams{
		AccountID: accountId,
		Amount:    amount,
	}
	res, err := s.repo.Withdraw(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Uint64("account_id", accountId).Msg("failed to create withdraw transaction")
		return model.TransactionDetails{}, err
	}

	_, err = s.accountService.Withdraw(ctx, accountId, int64(amount*100), res.Transaction.ID)
	if err != nil {
		s.logger.Error().Err(err).Uint64("account_id", accountId).Uint64("transaction_id", res.Transaction.ID).Msg("failed to create withdraw operation")
		updateErr := s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusFailed)
		if updateErr != nil {
			s.logger.Error().Err(updateErr).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update withdraw operation")
		}
		return model.TransactionDetails{}, updateErr
	}
	err = s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusCompleted)
	if err != nil {
		s.logger.Error().Err(err).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update withdraw operation")
		return model.TransactionDetails{}, err
	}
	res.Transaction.Status = model.TransactionStatusCompleted
	return res, nil
}

func (s *TransactionService) Transfer(ctx context.Context, userId, recipient uint64, amount float64) (model.TransactionDetails, error) {
	params := model.TransferParams{
		UserID:    userId,
		Recipient: recipient,
		Amount:    amount,
	}
	res, err := s.repo.Transfer(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Uint64("user_id", userId).Msg("failed to create transfer transaction")
		return model.TransactionDetails{}, err
	}

	_, err = s.accountService.Transfer(ctx, userId, recipient, int64(amount*100), res.Transaction.ID)
	if err != nil {
		s.logger.Error().Err(err).Uint64("user_id", userId).Msg("failed to create transfer operation")
		updateErr := s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusFailed)
		if updateErr != nil {
			s.logger.Error().Err(updateErr).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update transaction status")
		}
		return model.TransactionDetails{}, updateErr
	}
	err = s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusCompleted)
	if err != nil {
		s.logger.Error().Err(err).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update transaction status")
		return model.TransactionDetails{}, err
	}
	res.Transaction.Status = model.TransactionStatusCompleted
	return res, nil

}
