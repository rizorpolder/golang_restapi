package server

import (
	"context"
	"github.com/rs/zerolog"
	authpb "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/auth/model"
	"golang_restapi/internal/auth/service"
)

type Server struct {
	authpb.UnimplementedAuthServer
	authService service.AuthService
	logger      *zerolog.Logger
}

func New(authService service.AuthService, logger *zerolog.Logger) *Server {
	return &Server{authService: authService, logger: logger}
}

type AuthService interface {
	Login(context.Context, authpb.LoginRequest) (authpb.TokenPair, error)
	Logout(context.Context, string) error
	Register(context.Context, authpb.RegisterRequest) (model.User, error)
	Refresh(context.Context, string) (authpb.TokenPair, error)
	ValidateToken(string) (uint64, error)
	Delete(context.Context, uint64) error
}
