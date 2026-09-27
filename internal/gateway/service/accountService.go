package service

import (
	"context"
	"fmt"
	"golang_restapi/internal/gateway/account/mapper"
	"golang_restapi/internal/gateway/model"
)

// методы управления пользователями

func (s *GatewayService) CreateUser(ctx context.Context, newUser model.CreateUser) error {
	user := mapper.CreateUserToUser(newUser)
	_, err := s.accountService.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed create user: %w", err)
	}
	return nil
}

func (s *GatewayService) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	return s.accountService.GetUsers(ctx, limit, offset)
}

func (s *GatewayService) GetUser(ctx context.Context, id uint64) (model.User, error) {
	return s.accountService.GetUser(ctx, id)
}
func (s *GatewayService) DeleteUser(ctx context.Context, id uint64) error {
	err := s.authService.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed delete user: %w", err)
	}
	return s.accountService.DeleteUser(ctx, id)
}
func (s *GatewayService) UpdateUser(ctx context.Context, userId uint64, user model.User) error {
	return s.accountService.UpdateUser(ctx, userId, user)
}
