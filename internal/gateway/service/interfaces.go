package service

import (
	"context"
	"github.com/rs/zerolog"
	"golang_restapi/internal/gateway/model"
)

type GatewayService struct {
	logger         *zerolog.Logger
	accountService AccountService
	authService    AuthService
}

type AccountService interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUser(ctx context.Context, id uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	UpdateUser(ctx context.Context, id uint64, user model.User) error
}

type AuthService interface {
	Login(ctx context.Context, username string, password string) (model.TokenPair, error)
	Logout(context.Context, string) error
	Register(context.Context, uint64, string, string, string) error
	Verify(ctx context.Context, accessToken string) (uint64, error)
	RefreshToken(ctx context.Context, refreshToken string) (model.TokenPair, error)
	DeleteUser(ctx context.Context, id uint64) error
}
