package service

import (
	"context"
	"github.com/rs/zerolog"
	"golang_restapi/internal/account/model"
	"golang_restapi/internal/account/repository"
	"time"
)

type AccountService struct {
	repo   repository.Repository
	logger *zerolog.Logger
}

func New(repo repository.Repository, logger *zerolog.Logger) *AccountService {
	return &AccountService{repo: repo, logger: logger}
}

type Repository interface {
	CreateUser(context.Context, model.User) error
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error
}

func (s *AccountService) CreateUser(ctx context.Context, newUser model.CreateUser) error {
	user := model.User{
		Login:      newUser.Login,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		MiddleName: newUser.MiddleName,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.repo.CreateUser(ctx, user)
}
func (s *AccountService) GetUser(ctx context.Context, id uint64) (model.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *AccountService) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}
func (s *AccountService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.DeleteUser(ctx, id)
}
func (s *AccountService) UpdateUser(ctx context.Context, userId uint64, user model.UpdateUser) error {
	return s.repo.UpdateUser(ctx, userId, user)
}
