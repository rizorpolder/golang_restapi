package auth

import (
	"context"
	"fmt"
	authpb "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/gateway/auth/mapper"
	"golang_restapi/internal/gateway/model"
)

type Service struct {
	client authpb.AuthClient
}

func NewService(client authpb.AuthClient) *Service {
	return &Service{client: client}
}

func (s *Service) Login(ctx context.Context, loginOrEmail string, password string) (model.TokenPair, error) {
	res, err := s.client.Login(ctx, &authpb.LoginRequest{
		LoginOrEmail: loginOrEmail,
		Password:     password,
	})
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("login failed: %w", err)
	}
	return mapper.PbToTokenPair(res.TokenPair), nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	res, err := s.client.Refresh(ctx, &authpb.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("refresh failed: %w", err)
	}
	return mapper.PbToTokenPair(res.TokenPair), nil
}

func (s *Service) Verify(ctx context.Context, accessToken string) (uint64, error) {
	res, err := s.client.ValidateToken(ctx, &authpb.ValidateRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return 0, fmt.Errorf("failed verify access token: %w", err)
	}
	return res.UserId, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, er := s.client.Logout(ctx, &authpb.LogoutRequest{
		RefreshToken: refreshToken,
	})
	if er != nil {
		return fmt.Errorf("logout failed: %w", er)
	}
	return nil
}
