package account

import (
	"context"
	"fmt"
	accountpb "golang_restapi/contracts/account/go"
	pagination "golang_restapi/contracts/pagination/go"
	"golang_restapi/internal/gateway/account/mapper"
	"golang_restapi/internal/gateway/model"
)

type Service struct {
	client accountpb.AccountClient
}

func NewService(client accountpb.AccountClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	res, err := s.client.CreateUser(ctx, &accountpb.CreateUserRequest{
		User: mapper.UserCreateToPb(user),
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed create user: %w", err)
	}
	return mapper.PbToUser(res.User), nil
}
func (s *Service) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	user, err := s.client.GetUser(ctx, &accountpb.GetUserRequest{UserId: userID})
	if err != nil {
		return model.User{}, fmt.Errorf("failed get user: %w", err)
	}
	return mapper.PbToUser(user.User), nil
}

func (s *Service) GetUsers(ctx context.Context, limit uint32, offset uint32) ([]model.User, error) {

	pagination := &pagination.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	users, err := s.client.GetUsers(ctx, &accountpb.GetUsersRequest{
		Pagination: pagination})

	if err != nil {
		return nil, fmt.Errorf("failed get users: %w", err)

	}
	return mapper.PbsToUsers(users.Users), nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uint64) error {
	_, err := s.client.DeleteUser(ctx, &accountpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("failed delete user: %w", err)
	}
	return nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uint64, user model.User) error {
	_, err := s.client.UpdateUser(ctx, &accountpb.UpdateUserRequest{
		UserId: userID,
		User:   mapper.UserUpdateToPb(user),
	})
	if err != nil {
		return fmt.Errorf("failed update user: %w", err)
	}
	return nil
}
