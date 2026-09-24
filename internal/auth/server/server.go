package server

import (
	"context"
	"github.com/rs/zerolog"
	authpb "golang_restapi/contracts/auth/go"
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
	Register(ctx context.Context) error
	Login(ctx context.Context) (string, error)
	Refresh(ctx context.Context)
	ValidateToken(ctx context.Context)
	Logout(ctx context.Context)
}
