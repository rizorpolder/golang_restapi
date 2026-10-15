package server

import (
	"context"
	"github.com/rs/zerolog"
	common "golang_restapi/contracts/common/go"
	transactionpb "golang_restapi/contracts/transaction/go"
	"golang_restapi/internal/transaction/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Server struct {
	transactionpb.UnimplementedTransactionServiceServer
	transactionService TransactionService
	logger             *zerolog.Logger
}

func New(transactionService TransactionService, logger *zerolog.Logger) *Server {
	return &Server{
		transactionService: transactionService,
		logger:             logger,
	}
}

type TransactionService interface {
	GetTransactionsWithDetails(context.Context, model.GetTransactionParams) ([]model.TransactionDetails, error)
	Deposit(context.Context, uint64, float64) (model.TransactionDetails, error)
	Withdraw(context.Context, uint64, float64) (model.TransactionDetails, error)
	Transfer(context.Context, uint64, uint64, float64) (model.TransactionDetails, error)
}

func (s *Server) Deposit(ctx context.Context, req *transactionpb.DepositRequest) (*common.EmptyResponse, error) {
	_, err := s.transactionService.Deposit(ctx, req.UserId, float64(req.Amount))
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}
func (s *Server) Withdraw(ctx context.Context, req *transactionpb.WithdrawRequest) (*common.EmptyResponse, error) {
	_, err := s.transactionService.Withdraw(ctx, req.UserId, float64(req.Amount))
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}
func (s *Server) Transfer(ctx context.Context, req *transactionpb.TransferRequest) (*common.EmptyResponse, error) {
	_, err := s.transactionService.Transfer(ctx, req.UserId, req.Recipient, float64(req.Amount))
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}
func (s *Server) GetTransactions(ctx context.Context, req *transactionpb.GetTransactionsRequest) (*transactionpb.GetTransactionsResponse, error) {

	params := model.GetTransactionParams{
		UserID: req.UserId,
		Type:   req.Type,
		Status: req.Status,
		Limit:  int(req.Pagination.Limit),
		Offset: int(req.Pagination.Offset),
	}

	if req.DateFrom != nil {
		t := req.DateFrom.AsTime()
		params.DateFrom = &t
	}

	if req.DateTo != nil {
		t := req.DateTo.AsTime()
		params.DateTo = &t
	}
	
	details, err := s.transactionService.GetTransactionsWithDetails(ctx, params)
	if err != nil {
		return &transactionpb.GetTransactionsResponse{}, err
	}
	transactions := make([]*transactionpb.Transaction, len(details))
	for index, detail := range details {
		//Конвертировать транзакции в pb и отправить
	}
	&transactionpb.GetTransactionsResponse{}
}
