package account

import (
	"context"
	"errors"
	"fmt"
	"golang_restapi/internal/account/model"
	"golang_restapi/internal/repository/account/mapper"
	repomodel "golang_restapi/internal/repository/account/model"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (repo *Repository) CreateUser(ctx context.Context, user model.User) error {
	userRepos := mapper.UserToRepoUser(user)
	res := repo.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&userRepos)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to save user")
		return fmt.Errorf("failed to save user %w", res.Error)
	}
	return nil
}

func (repo *Repository) GetUser(ctx context.Context, id uint64) (model.User, error) {
	var user repomodel.User
	res := repo.db.WithContext(ctx).Model(&user).Where("id = ?", id).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, fmt.Errorf("User not found") //404
	} else if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get user id")
	}
	return mapper.RepoUserToUser(user), nil
}

func (repo *Repository) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var repoUsers []repomodel.User
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Offset(offset).
		Limit(limit).
		Find(&repoUsers)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("users not found")
	} else if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get repoUsers")
	}

	return mapper.RepoUsersToUsers(repoUsers), nil

}

func (repo *Repository) DeleteUser(ctx context.Context, id uint64) error {
	res := repo.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&repomodel.User{})
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to delete user")
		return fmt.Errorf("failed to delete user %w", res.Error)
	}
	return nil
}

func (repo *Repository) UpdateUser(ctx context.Context, user model.User) error {
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", user.ID).
		Updates(user)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to update user")
		return fmt.Errorf("failed to update user %w", res.Error)
	}
	return nil
}
