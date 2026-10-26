package server

import (
	"context"
	"github.com/rs/zerolog"
	accountpb "golang_restapi/contracts/account/go"
	commonpb "golang_restapi/contracts/common/go"
	"golang_restapi/internal/account/mapper"
	"golang_restapi/internal/account/model"
)

type Server struct {
	accountpb.UnimplementedAccountServer
	accountService AccountService
	logger         *zerolog.Logger
}

func New(accountService AccountService, logger *zerolog.Logger) *Server {
	return &Server{accountService: accountService, logger: logger}
}

type AccountService interface {
	CreateUser(context.Context, model.CreateUser) (model.User, error)
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error
	GetBalance(context.Context, uint64) (model.GetBalanceResponse, error)
	Deposit(ctx context.Context, userId uint64, amount int64) (model.UpdateBalanceResponse, error)
	Withdraw(ctx context.Context, userId uint64, amount int64) (model.UpdateBalanceResponse, error)
	Transfer(ctx context.Context, fromUser, ToUser uint64, amount int64) (model.GetBalanceResponse, model.GetBalanceResponse, error)
}

func (s *Server) CreateUser(ctx context.Context, req *accountpb.CreateUserRequest) (*accountpb.CreateUserResponse, error) {
	user := mapper.PbToUserCreate(req.GetUser())
	result, err := s.accountService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &accountpb.CreateUserResponse{User: mapper.UserToPb(result)}, nil
}

func (s *Server) GetUser(ctx context.Context, req *accountpb.GetUserRequest) (*accountpb.GetUserResponse, error) {
	res, err := s.accountService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &accountpb.GetUserResponse{User: mapper.UserToPb(res)}, nil
}

func (s *Server) GetUsers(ctx context.Context, req *accountpb.GetUsersRequest) (*accountpb.GetUsersResponse, error) {
	res, err := s.accountService.GetUsers(ctx, int(req.Pagination.GetLimit()), int(req.GetPagination().GetOffset()))
	if err != nil {
		return nil, err
	}
	return &accountpb.GetUsersResponse{Users: mapper.UsersToPbs(res)}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *accountpb.UpdateUserRequest) (*commonpb.EmptyResponse, error) {
	user := mapper.PbToUserUpdate(req.User)
	err := s.accountService.UpdateUser(ctx, uint64(req.GetUserId()), user)
	if err != nil {
		return nil, err
	}
	return &commonpb.EmptyResponse{}, nil
}
func (s *Server) DeleteUser(ctx context.Context, req *accountpb.DeleteUserRequest) (*accountpb.DeleteUserResponse, error) {
	err := s.accountService.DeleteUser(ctx, uint64(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &accountpb.DeleteUserResponse{Result: true}, nil
}

func (s *Server) Deposit(ctx context.Context, req *accountpb.DepositRequest) (*accountpb.DepositResponse, error) {

	response, err := s.accountService.Deposit(ctx, req.UserId, req.Amount)
	if err != nil {
		return &accountpb.DepositResponse{
			Status:  "failed",
			Balance: 0,
		}, err
	}
	return &accountpb.DepositResponse{
		Status:  "completed",
		Balance: int64(response.NewBalance),
	}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *accountpb.WithdrawRequest) (*accountpb.WithdrawResponse, error) {
	response, err := s.accountService.Withdraw(ctx, req.UserId, req.Amount)
	if err != nil {
		return &accountpb.WithdrawResponse{
			Status:  "failed",
			Balance: 0,
		}, nil
	}
	return &accountpb.WithdrawResponse{
		Status:  "completed",
		Balance: int64(response.NewBalance),
	}, nil
}

func (s *Server) Transfer(ctx context.Context, req *accountpb.TransferRequest) (*accountpb.TransferResponse, error) {
	fromUserBalance, toUserBalance, err := s.accountService.Transfer(ctx, req.UserId, req.RecipientId, req.Amount)
	if err != nil {
		return &accountpb.TransferResponse{
			Status:           "failed",
			UserBalance:      0,
			RecipientBalance: 0,
		}, err
	}
	return &accountpb.TransferResponse{
		Status:           "completed",
		UserBalance:      int64(fromUserBalance.Balance),
		RecipientBalance: int64(toUserBalance.Balance),
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.GetBalanceResponse, error) {
	res, err := s.accountService.GetBalance(ctx, req.UserId)
	if err != nil {
		return &accountpb.GetBalanceResponse{
			Balance: 0,
		}, err
	}
	return &accountpb.GetBalanceResponse{
		Balance: int64(res.Balance),
	}, nil
}
