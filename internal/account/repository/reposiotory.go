package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"golang_restapi/internal/account/model"
	"golang_restapi/internal/account/repository/mapper"
	"golang_restapi/internal/account/repository/model"
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

func (repo *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	var result repomodel.User
	userRepos := mapper.UserToRepoUser(user)
	res := repo.db.WithContext(ctx).Model(&result).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&userRepos)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to save user")
		return model.User{}, fmt.Errorf("failed to save user %w", res.Error)
	}
	return mapper.RepoUserToUser(result), nil
}

func (repo *Repository) GetUser(ctx context.Context, id uint64) (model.User, error) {
	var user repomodel.User
	res := repo.db.WithContext(ctx).Model(&user).Where("id = ? AND is_deleted = ?", id, false).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, fmt.Errorf("user not found") //404
	} else if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get user id")
	}
	return mapper.RepoUserToUser(user), nil
}

func (repo *Repository) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var repoUsers []repomodel.User
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Find(&repoUsers)

	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get repoUsers")
		return nil, res.Error
	}

	if len(repoUsers) == 0 {
		return nil, fmt.Errorf("users not found")
	}

	return mapper.RepoUsersToUsers(repoUsers), nil

}

func (repo *Repository) DeleteUser(ctx context.Context, id uint64) error {
	res := repo.db.WithContext(ctx).
		Where("id = ?", id).
		Update("is_deleted", true)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to delete user")
		return fmt.Errorf("failed to delete user %w", res.Error)
	}
	return nil
}

func (repo *Repository) UpdateUser(ctx context.Context, userId uint64, user model.UpdateUser) error {
	repoUser := mapper.UpdateUserToRepoUser(user)
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ? AND is_deleted = ?", userId, false).
		Updates(repoUser)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to update user")
		return fmt.Errorf("failed to update user %w", res.Error)
	}
	return nil
}

func (repo *Repository) GetBalance(ctx context.Context, userID uint64) (float32, error) {
	var repUser repomodel.User
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Select("balance").
		Where("id = ? AND is_deleted = ?", userID, false).
		First(&repUser)
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get balance")
		return 0, fmt.Errorf("failed to get balance %w", res.Error)
	}
	return repUser.Balance, nil
}

func (repo *Repository) UpdateBalance(ctx context.Context, userId uint64, amount float32, operationType model.BalanceOperationType) (float32, float32, error) {
	var user repomodel.User
	res := repo.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ? AND is_deleted = ?", userId, false).
		Select("balance").
		First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, 0, fmt.Errorf("user not found")
	} else if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to get user")
		return 0, 0, fmt.Errorf("failed to get balance %w", res.Error)
	}

	oldBalance := user.Balance
	var newBalance float32
	switch operationType {
	case model.BalanceOperationDeposit:
		newBalance = oldBalance + amount
	case model.BalanceOperationCredit:
		newBalance = oldBalance - amount
		if newBalance < 0 {
			return 0, 0, fmt.Errorf("insufficient balance")
		}
	default:
		return 0, 0, fmt.Errorf("invalid operation type")
	}
	if res.Error != nil {
		repo.logger.Err(res.Error).Msg("failed to update balance")
		return 0, 0, fmt.Errorf("failed to update balance %w", res.Error)
	}
	return oldBalance, newBalance, nil
}

func (repo *Repository) TransferBalance(ctx context.Context, fromUserID, toUserID uint64, amount float32) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return fmt.Errorf("cannot transfer to the same user")
	}

	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var toUser repomodel.User
		res := tx.Model(&repomodel.User{}).
			Where("id = ? AND is_deleted = ?", toUserID, false).
			First(&toUser)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sender user not found")
		} else if res.Error != nil {
			return fmt.Errorf("failed to find user: %w", res.Error)
		}

		var fromUser repomodel.User
		res = tx.Model(&repomodel.User{}).
			Where("id = ? AND is_deleted = ?", fromUserID, false).
			First(&fromUser)

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sender user not found")
		} else if res.Error != nil {
			return fmt.Errorf("failed to find user: %w", res.Error)
		}

		if fromUser.Balance < amount {
			return fmt.Errorf("insufficient balance: current balance %.2f, transfer amount %.2f", fromUser.Balance, amount)
		}

		res = tx.Model(&repomodel.User{}).
			Where("id = ? AND is_deleted = ?", fromUserID, false).
			Update("balance", fromUser.Balance-amount)
		if res.Error != nil {
			return fmt.Errorf("failed to deduct sender balance %w", res.Error)
		}

		res = tx.Model(&repomodel.User{}).
			Where("id = ? AND is_deleted = ?", toUserID, false).
			Update("balance", toUser.Balance+amount)
		if res.Error != nil {
			return fmt.Errorf("failed to add to recipient balance %w", res.Error)
		}
		return nil
	})
}
