package server

import (
	"context"
	"github.com/rs/zerolog"
	accountpb "golang_restapi/contracts/account/go"
	common "golang_restapi/contracts/common/go"
	gatewaypb "golang_restapi/contracts/gateway/go"
	accountMapper "golang_restapi/internal/gateway/account/mapper"
	authMapper "golang_restapi/internal/gateway/auth/mapper"
	"golang_restapi/internal/gateway/model"
	"golang_restapi/internal/gateway/service"
)

type Server struct {
	gatewaypb.UnimplementedGatewayServer
	gatewayService service.GatewayService
	logger         *zerolog.Logger
}

func New(gatewayService service.GatewayService, logger *zerolog.Logger) *Server {
	return &Server{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

type GatewayService interface {
	//Auth
	Register(ctx context.Context, newUser model.User, password string) (model.User, model.TokenPair, error)
	Login(ctx context.Context, loginOrEmail string, password string) (model.User, model.TokenPair, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error)
	ValidateToken(ctx context.Context, accessToken string) (uint64, bool, error)
	//Account
	CreateUser(ctx context.Context, newUser model.CreateUser) (model.User, error)
	GetUsers(ctx context.Context, limit uint32, offset uint32) ([]model.User, error)
	GetUser(ctx context.Context, id uint64) (model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	UpdateUser(ctx context.Context, userId uint64, user model.User) error
}

func (s *Server) Register(ctx context.Context, req *gatewaypb.RegisterRequest) (*gatewaypb.RegisterResponse, error) {

	user := accountMapper.PbToUser(req.User)
	createUser, tokens, err := s.gatewayService.Register(ctx, user, req.Password)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.RegisterResponse{
		User:      accountMapper.UserToPb(createUser),
		TokenPair: authMapper.TokenPairToPb(tokens),
	}, nil
}

func (s *Server) Login(ctx context.Context, req *gatewaypb.LoginRequest) (*gatewaypb.LoginResponse, error) {
	user, tokens, err := s.gatewayService.Login(ctx, req.LoginOrEmail, req.Password)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.LoginResponse{
		User:      accountMapper.UserToPb(user),
		TokenPair: authMapper.TokenPairToPb(tokens),
	}, nil
}
func (s *Server) Refresh(ctx context.Context, req *gatewaypb.RefreshRequest) (*gatewaypb.RefreshResponse, error) {
	tokenPair, err := s.gatewayService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.RefreshResponse{TokenPair: authMapper.TokenPairToPb(tokenPair)}, nil
}
func (s *Server) Logout(ctx context.Context, req *gatewaypb.LogoutRequest) (*common.EmptyResponse, error) {
	err := s.gatewayService.Logout(ctx, req.RefreshToken)
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *gatewaypb.ValidateTokenRequest) (*gatewaypb.ValidateTokenResponse, error) {
	userId, result, err := s.gatewayService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.ValidateTokenResponse{
		UserId:  userId,
		IsValid: result,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *gatewaypb.CreateUserRequest) (*gatewaypb.CreateUserResponse, error) {
	user := accountMapper.PbToCreateUser(req.User)
	created, err := s.gatewayService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.CreateUserResponse{
		User: accountMapper.UserToPb(created),
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *gatewaypb.GetUserRequest) (*gatewaypb.GetUserResponse, error) {
	user, err := s.gatewayService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &gatewaypb.GetUserResponse{
		User: accountMapper.UserToPb(user),
	}, nil
}

func (s *Server) GetCurrentUser(ctx context.Context, request *common.EmptyRequest) (*gatewaypb.GetCurrentUserResponse, error) {
	return nil, nil
}

func (s *Server) GetUsers(ctx context.Context, req *gatewaypb.GetUsersRequest) (*gatewaypb.GetUsersResponse, error) {
	users, err := s.gatewayService.GetUsers(ctx, req.Pagination.Limit, req.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*accountpb.User, len(users))
	for i := range users {
		pbUsers[i] = accountMapper.UserToPb(users[i])
	}
	return &gatewaypb.GetUsersResponse{Users: pbUsers}, nil
}

func (s *Server) UpdateCurrentUser(ctx context.Context, req *gatewaypb.UpdateCurrentUserRequest) (*common.EmptyResponse, error) {
	return nil, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *gatewaypb.UpdateUserRequest) (*common.EmptyResponse, error) {
	user := accountMapper.PbToUser(req.User)
	err := s.gatewayService.UpdateUser(ctx, req.UserId, user)
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *gatewaypb.DeleteUserRequest) (*common.EmptyResponse, error) {
	err := s.gatewayService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return &common.EmptyResponse{}, err
	}
	return &common.EmptyResponse{}, nil
}
