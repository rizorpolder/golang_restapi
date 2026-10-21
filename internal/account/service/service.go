package service

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_restapi/internal/account/config"
	"golang_restapi/internal/account/model"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type AccountService struct {
	repo           Repository
	cfg            *config.Config
	logger         *zerolog.Logger
	kafkaPublisher KafkaPublisher
}

func New(repo Repository, kafkaPublisher KafkaPublisher, cfg *config.Config, logger *zerolog.Logger) *AccountService {
	return &AccountService{repo: repo, kafkaPublisher: kafkaPublisher, cfg: cfg, logger: logger}
}

type TransactionRequest struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	Amount      int64  `json:"amount"`
	OperationID uint64 `json:"operation_id"`
	RecipientID uint64 `json:"recipient_id"`
}

type TransactionResponse struct {
	RequestType string                 `json:"request_type"`
	UserID      uint64                 `json:"user_id"`
	OperationID uint64                 `json:"operation_id"`
	Result      map[string]interface{} `json:"result"`
}

type Repository interface {
	CreateUser(context.Context, model.User) (model.User, error)
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error

	GetBalance(context.Context, uint64) (float32, error)
	UpdateBalance(context.Context, uint64, float32, model.BalanceOperationType) (float32, float32, error)
	TransferBalance(context.Context, uint64, uint64, float32) error
}

type KafkaPublisher interface {
	Publish(ctx context.Context, topic string, key string, data interface{}) error
}

func (s *AccountService) CreateUser(ctx context.Context, newUser model.CreateUser) (model.User, error) {
	user := model.User{
		Login:      newUser.Login,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		MiddleName: newUser.MiddleName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.repo.CreateUser(ctx, user)
}
func (s *AccountService) GetUser(ctx context.Context, id uint64) (model.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *AccountService) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}
func (s *AccountService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.DeleteUser(ctx, id)
}
func (s *AccountService) UpdateUser(ctx context.Context, userId uint64, user model.UpdateUser) error {
	return s.repo.UpdateUser(ctx, userId, user)
}

func (s *AccountService) GetBalance(ctx context.Context, userId uint64) (model.GetBalanceResponse, error) {
	balance, err := s.repo.GetBalance(ctx, userId)
	if err != nil {
		s.logger.Err(err).Uint64("userId", userId).Msg("failed to get user balance")
		return model.GetBalanceResponse{}, fmt.Errorf("failed to get user balance: %w", err)
	}
	return model.GetBalanceResponse{
		UserID:  userId,
		Balance: balance,
	}, nil
}
func (s *AccountService) UpdateBalance(ctx context.Context, req model.UpdateBalanceRequest) (model.UpdateBalanceResponse, error) {
	if req.Amount < 0 {
		return model.UpdateBalanceResponse{}, fmt.Errorf("amount must be greater than zero")
	}
	if req.Type == model.BalanceOperationCredit && req.Amount == 0 {
		return model.UpdateBalanceResponse{}, fmt.Errorf("amount must be positive for substract")
	}

	oldBalance, newBalance, err := s.repo.UpdateBalance(ctx, req.UserID, req.Amount, req.Type)
	if err != nil {
		s.logger.Err(err).Uint64("userId", req.UserID).
			Float32("amount", req.Amount).
			Str("operation", string(req.Type)).
			Msg("failed to update balance")
		return model.UpdateBalanceResponse{}, fmt.Errorf("failed to update balance: %w", err)
	}

	s.logger.Info().
		Uint64("userId", req.UserID).
		Float32("amount", req.Amount).
		Str("operation", string(req.Type)).
		Msg("updated balance successfully")
	return model.UpdateBalanceResponse{
		UserID:        req.UserID,
		OldBalance:    oldBalance,
		NewBalance:    newBalance,
		Amount:        req.Amount,
		OperationType: req.Type,
	}, nil
}

func (s *AccountService) TransferBalance(ctx context.Context, fromUserID, toUserID uint64, amount float32) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return fmt.Errorf("cannot transfer to same user")
	}
	err := s.repo.TransferBalance(ctx, fromUserID, toUserID, amount)
	if err != nil {
		s.logger.Err(err).Uint64("fromUserId", fromUserID).
			Uint64("toUserId", toUserID).
			Float32("amount", amount).
			Msg("failed to transfer balance")
		return err
	}
	s.logger.Info().Uint64("fromUserId", fromUserID).
		Uint64("toUserId", toUserID).
		Float32("amount", amount).
		Msg("transfer balance successfully")
	return nil
}

func (s *AccountService) Deposit(ctx context.Context, userID uint64, amount int64) (model.UpdateBalanceResponse, error) {
	updateReq := model.UpdateBalanceRequest{
		UserID: userID,
		Amount: float32(amount),
		Type:   model.BalanceOperationDeposit,
	}

	response, err := s.UpdateBalance(ctx, updateReq)
	if err != nil {
		s.logger.Err(err).
			Uint64("userID", userID).
			Int64("amount", amount).
			Msg("failed to deposit balance")
		return model.UpdateBalanceResponse{}, err
	}
	s.logger.Info().Uint64("userID", userID).
		Int64("amount", amount).
		Msg("deposit balance successfully")
	return response, nil
}

func (s *AccountService) Withdraw(ctx context.Context, userID uint64, amount int64) (model.UpdateBalanceResponse, error) {
	updateReq := model.UpdateBalanceRequest{
		UserID: userID,
		Amount: float32(amount),
		Type:   model.BalanceOperationCredit,
	}
	response, err := s.UpdateBalance(ctx, updateReq)
	if err != nil {
		s.logger.Err(err).
			Uint64("userID", userID).
			Int64("amount", amount).
			Msg("failed to withdraw balance")
		return model.UpdateBalanceResponse{}, err
	}
	s.logger.Info().Uint64("userID", userID).
		Int64("amount", amount).
		Float32("newBalance", response.NewBalance).
		Msg("withdraw balance successfully")
	return response, nil
}

func (s *AccountService) Transfer(ctx context.Context, fromUserID, toUserID uint64, amount int64) (model.GetBalanceResponse, model.GetBalanceResponse, error) {
	err := s.TransferBalance(ctx, fromUserID, toUserID, float32(amount))
	if err != nil {
		s.logger.Err(err).Uint64("fromUserId", fromUserID).
			Uint64("toUserId", toUserID).
			Int64("amount", amount).
			Msg("failed to transfer balance")
		return model.GetBalanceResponse{}, model.GetBalanceResponse{}, err
	}
	fromBalance, err := s.GetBalance(ctx, fromUserID)
	if err != nil {
		s.logger.Err(err).Uint64("fromUserId", fromUserID).Msg("failed to get from balance")
		return model.GetBalanceResponse{}, model.GetBalanceResponse{}, err
	}
	toBalance, err := s.GetBalance(ctx, toUserID)
	if err != nil {
		s.logger.Err(err).Uint64("toUserId", toUserID).Msg("failed to get to balance")
		return model.GetBalanceResponse{}, model.GetBalanceResponse{}, err
	}
	s.logger.Info().Uint64("fromUserId", fromUserID).
		Uint64("toUserId", toUserID).
		Int64("amount", amount).
		Msg("transfer balance successfully")
	return fromBalance, toBalance, nil
}

func (s *AccountService) HandleTransaction(ctx context.Context, topic string, key string, data []byte) error {
	s.logger.Info().Msg("handling transaction request")
	var res TransactionRequest
	var err error

	err = json.Unmarshal(data, &res)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	result := map[string]interface{}{
		"request_type": res.RequestType,
		"user_id":      res.UserID,
		"operation_id": res.OperationID,
		"result":       false,
	}
	switch res.RequestType {
	case "withdraw":
		_, err = s.Withdraw(ctx, res.UserID, res.Amount)
	case "deposit":
		_, err = s.Deposit(ctx, res.UserID, res.Amount)
	case "transfer":
		_, _, err = s.Transfer(ctx, res.UserID, res.RecipientID, res.Amount)
	default:
		err = fmt.Errorf("unknown transaction request type: %s", res.RequestType)
	}
	if err != nil {
		s.logger.Err(err).
			Uint64("userID", res.UserID).
			Uint64("operationID", res.OperationID).
			Msg("failed to handle transaction")
		s.kafkaPublisher.Publish(ctx, s.cfg.KafkaPublishTopic, strconv.FormatUint(res.UserID, 10), result)
		return fmt.Errorf("failed to handle transaction: %w", err)
	}
	result["result"] = true
	s.kafkaPublisher.Publish(ctx, s.cfg.KafkaPublishTopic, strconv.FormatUint(res.UserID, 10), result)
	return nil
}
