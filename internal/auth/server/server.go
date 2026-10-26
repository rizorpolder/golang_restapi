package server

import (
	"context"
	"github.com/rs/zerolog"
	authpb "golang_restapi/contracts/auth/go"
	common "golang_restapi/contracts/common/go"
	"golang_restapi/internal/auth/mapper"
	"golang_restapi/internal/auth/model"
)

type Server struct {
	authpb.UnimplementedAuthServer
	authService AuthService
	logger      *zerolog.Logger
}

func New(authService AuthService, logger *zerolog.Logger) *Server {
	return &Server{authService: authService, logger: logger}
}

type AuthService interface {
	Register(context.Context, model.Register) error
	Login(context.Context, model.Login) (model.TokenPair, error)
	Logout(context.Context, string) error
	Refresh(context.Context, string) (model.TokenPair, error)
	Validate(context.Context, string) (uint64, error)
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*common.EmptyResponse, error) {
	register := mapper.PbRegisterToRegister(req)
	if err := s.authService.Register(ctx, register); err != nil {
		return nil, err
	}
	return &common.EmptyResponse{}, nil
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.TokenPair, error) {
	login := mapper.PbLoginToLogin(req)

	tokenPair, err := s.authService.Login(ctx, login)
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToPbTokenPair(tokenPair), nil
}
func (s *Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.TokenPair, error) {
	tokenPair, err := s.authService.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return mapper.TokenPairToPbTokenPair(tokenPair), nil
}

func (s *Server) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	userId, err := s.authService.Validate(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateResponse{UserId: userId}, nil
}
func (s *Server) Logout(ctx context.Context, req *authpb.RefreshRequest) (*common.EmptyResponse, error) {
	if err := s.authService.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}
	return &common.EmptyResponse{}, nil
}
