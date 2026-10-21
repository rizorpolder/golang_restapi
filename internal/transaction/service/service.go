package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_restapi/internal/transaction/config"
	"golang_restapi/internal/transaction/model"
	"time"

	"github.com/rs/zerolog"
)

type TransactionService struct {
	cfg    *config.Config
	repo   Repository
	logger *zerolog.Logger

	accountService AccountService
	kafkaPublisher KafkaPublisher
}

func New(repo Repository, accountService AccountService, kafkaPublisher KafkaPublisher, cfg *config.Config, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{
		repo:           repo,
		cfg:            cfg,
		accountService: accountService,
		kafkaPublisher: kafkaPublisher,
		logger:         logger,
	}
}

type AccountResponse struct {
	RequestType string `json:"request_type"`
	UserId      uint64 `json:"user_id"`
	OperationId uint64 `json:"operation_id"`
	Result      bool   `json:"result"`
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

type KafkaPublisher interface {
	Publish(ctx context.Context, topic string, key string, data interface{}) error
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

	request := map[string]interface{}{
		"request_type": "deposit",
		"user_id":      userId,
		"amount":       int64(amount * 100),
		"operation_id": res.Transaction.ID,
		"timestamp":    time.Now().UTC(),
	}

	return s.tryPublishToKafka(ctx, res, userId, request)
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

	request := map[string]interface{}{
		"request_type": "withdraw",
		"user_id":      accountId,
		"amount":       int64(amount * 100),
		"operation_id": res.Transaction.ID,
		"timestamp":    time.Now().UTC(),
	}

	return s.tryPublishToKafka(ctx, res, accountId, request)
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

	request := map[string]interface{}{
		"request_type": "transfer",
		"user_id":      userId,
		"amount":       int64(amount * 100),
		"operation_id": res.Transaction.ID,
		"recipient_id": recipient,
		"timestamp":    time.Now().UTC(),
	}

	return s.tryPublishToKafka(ctx, res, userId, request)

}

func (s *TransactionService) tryPublishToKafka(ctx context.Context, res model.TransactionDetails, userId uint64, request map[string]interface{}) (model.TransactionDetails, error) {
	err := s.kafkaPublisher.Publish(ctx, s.cfg.KafkaPublishTopic, fmt.Sprintf("%d", userId), request)
	if err != nil {
		s.logger.Error().Uint64("user_id", userId).Uint64("transaction_id", res.Transaction.ID).Msg("failed to publish transaction")
		updateErr := s.repo.UpdateTransactionStatus(ctx, res.Transaction.ID, model.TransactionStatusFailed)
		if updateErr != nil {
			s.logger.Error().Err(updateErr).Uint64("transaction_id", res.Transaction.ID).Msg("failed to update transaction status")
		}
		return model.TransactionDetails{}, updateErr
	}

	s.logger.Info().Uint64("user_id", userId).Uint64("transaction_id", res.Transaction.ID).Msgf("%v transaction success", res.Transaction.Type)
	return res, nil
}

func (s *TransactionService) HandleAccountResponse(ctx context.Context, topic string, key string, data []byte) error {
	var response AccountResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if response.RequestType == "" || response.UserId == 0 || response.OperationId == 0 {
		return fmt.Errorf("missing or invalid response")
	}
	s.logger.Info().
		Str("request_type", response.RequestType).
		Uint64("user_id", response.UserId).
		Uint64("operation_id", response.OperationId).
		Interface("result", response.Result).
		Msg("received account response via kafka")

	var status model.TransactionStatus
	var logMsg string
	if response.Result {
		status = model.TransactionStatusCompleted
		logMsg = "transaction success"
	} else {
		status = model.TransactionStatusFailed
		logMsg = "transaction failed"
	}

	err := s.repo.UpdateTransactionStatus(ctx, response.OperationId, status)
	if err != nil {
		s.logger.Error().Err(err).Uint64("operation_id", response.OperationId).Msgf("failed to update transaction status to %s", status)
		return fmt.Errorf("failed to update transaction status: to %s:%w", status, err)
	}
	s.logger.Info().
		Uint64("operation_id", response.OperationId).
		Msg(logMsg)
	return nil
}
