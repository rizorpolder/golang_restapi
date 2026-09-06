package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"golang_restapi/internal/auth/model"
	"golang_restapi/internal/auth/repository/mapper"
	repomodel "golang_restapi/internal/auth/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
	l  *zerolog.Logger
}

func NewRepository(db *gorm.DB, l *zerolog.Logger) *Repository {
	return &Repository{
		db: db,
		l:  l,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	userRepo := mapper.UserToRepoUser(user)
	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&userRepo)

	if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to create user")
		return fmt.Errorf("failed to create user: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userId uint64) (model.User, error) {
	var user repomodel.User
	res := r.db.WithContext(ctx).Model(&repomodel.User{}).
		Where("id = ?", userId).
		First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, errors.New("user not found")
	} else if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to get user")
		return model.User{}, fmt.Errorf("failed to get user: %w", res.Error)
	}
	return mapper.RepoUserToUser(user), nil
}

func (r *Repository) GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (model.User, error) {
	var user repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("login = ? OR email = ?", loginOrEmail, loginOrEmail).
		First(&user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, errors.New("user not found")
	} else if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to get user")
		return model.User{}, res.Error
	}

	return mapper.RepoUserToUser(user), nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error {
	m := mapper.TokenToRepoToken(refreshToken)
	res := r.db.WithContext(ctx).Create(&m)
	if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to save refresh token")
		return fmt.Errorf("failed to save refresh token: %w", res.Error)
	}
	return nil
}
func (r *Repository) GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	var repoToken repomodel.RefreshToken

	res := r.db.WithContext(ctx).Model(&repomodel.RefreshToken{}).
		Where("token = ?", token).
		First(&repoToken)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.RefreshToken{}, errors.New("refresh token not found")
	} else if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to get refresh token")
		return model.RefreshToken{}, res.Error
	}
	return mapper.RepoTokenToToken(repoToken), nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, token string) error {
	res := r.db.WithContext(ctx).Model(&repomodel.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", gorm.Expr("NOW()"))
	if res.Error != nil {
		r.l.Err(res.Error).Msg("Failed to revoke refresh token")
		return fmt.Errorf("failed to revoke refresh token: %w", res.Error)
	}
	return nil
}
