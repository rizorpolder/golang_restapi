package server

import (
	"context"
	"github.com/rs/zerolog"
	accountpb "golang_restapi/contracts/account/go"
	"golang_restapi/internal/account/mapper"
	"golang_restapi/internal/account/model"
	"golang_restapi/internal/account/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	accountpb.UnimplementedAccountServer
	accountService service.AccountService
	logger         *zerolog.Logger
}

func New(accountService service.AccountService, logger *zerolog.Logger) *Server {
	return &Server{accountService: accountService, logger: logger}
}

type AccountService interface {
	CreateUser(context.Context, model.User) error
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error
}

func (s *Server) CreateUser(ctx context.Context, req *accountpb.CreateUserRequest) (*emptypb.Empty, error) {
	user := mapper.PbToUserCreate(req.GetUser())
	if err := s.accountService.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
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

func (s *Server) UpdateUser(ctx context.Context, req *accountpb.UpdateUserRequest) (*emptypb.Empty, error) {
	user := mapper.PbToUserUpdate(req.User)
	err := s.accountService.UpdateUser(ctx, uint64(req.GetUserId()), user)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *Server) DeleteUser(ctx context.Context, req *accountpb.DeleteUserRequest) (*accountpb.DeleteUserResponse, error) {
	err := s.accountService.DeleteUser(ctx, uint64(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &accountpb.DeleteUserResponse{Result: true}, nil
}
