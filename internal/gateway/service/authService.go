package service

import (
	"context"
	"fmt"
	"golang_restapi/internal/gateway/model"
)

// методы аутентификации

func (s *GatewayService) Register(ctx context.Context, newUser model.User, password string) (model.User, model.TokenPair, error) {
	user, err := s.accountService.CreateUser(ctx, newUser)
	if err != nil {
		return model.User{}, model.TokenPair{}, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.authService.Register(ctx, user.ID, newUser.Login, newUser.Email, password); err != nil {
		return model.User{}, model.TokenPair{}, fmt.Errorf("failed to register user: %w", err)
	}
	token, err := s.authService.Login(ctx, newUser.Login, password)
	if err != nil {
		return model.User{}, model.TokenPair{}, fmt.Errorf("failed to login: %w", err)

	}
	return user, token, nil
}

func (s *GatewayService) Login(ctx context.Context, loginOrEmail string, password string) (model.User, model.TokenPair, error) {
	token, err := s.authService.Login(ctx, loginOrEmail, password)
	if err != nil {
		return model.User{}, model.TokenPair{}, fmt.Errorf("failed to login: %w", err)
	}
	user := model.User{}
	return user, token, nil
}

func (s *GatewayService) Logout(ctx context.Context, token string) error {
	if err := s.authService.Logout(ctx, token); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (s *GatewayService) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	return s.authService.RefreshToken(ctx, refreshToken)
}

func (s *GatewayService) ValidateToken(ctx context.Context, accessToken string) (uint64, bool, error) {
	userId, err := s.authService.Verify(ctx, accessToken)
	if err != nil {
		return 0, false, fmt.Errorf("failed to verify token: %w", err)
	}
	return userId, true, nil
}
