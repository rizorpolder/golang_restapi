package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	pbmodel "golang_restapi/contracts/auth/go"
	"golang_restapi/internal/auth/model"
	"golang_restapi/internal/auth/repository"
	"golang_restapi/internal/config"
	"time"
)

type AuthService struct {
	repository *repository.Repository
	cfg        config.Config
	logger     *zerolog.Logger
}

type Repository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	GetUserByLoginOrEmail(ctx context.Context, login string, email string) (model.User, error)
	SaveRefreshToken(ctx context.Context, token string, refreshToken string) (model.User, error)
	GetRefreshToken(ctx context.Context, token string) (model.User, error)
	RevokeRefreshToken(ctx context.Context, token string, refreshToken string) error
}

func New(repo *repository.Repository, cfg config.Config, logger *zerolog.Logger) *AuthService {
	return &AuthService{repository: repo, cfg: cfg, logger: logger}
}

func (service *AuthService) Register(ctx context.Context, req pbmodel.RegisterRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user := model.User{
		Login:        req.Login,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return service.repository.CreateUser(ctx, user)
}

func (service *AuthService) Login(ctx context.Context, req pbmodel.LoginRequest) (pbmodel.TokenPair, error) {
	user, err := service.repository.GetUserByLoginOrEmail(ctx, req.LoginOrEmail)
	if err != nil {
		return pbmodel.TokenPair{}, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return pbmodel.TokenPair{}, fmt.Errorf("invalid credentials")
	}

	return service.issueTokens(ctx, user.ID)
}

func (service *AuthService) Refresh(ctx context.Context, refreshToken string) (pbmodel.TokenPair, error) {
	rt, err := service.repository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return pbmodel.TokenPair{}, fmt.Errorf("invalid refresh token")
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return pbmodel.TokenPair{}, fmt.Errorf("refresh token is expired or revoked")
	}
	return service.issueTokens(ctx, rt.ID)
}

func (service *AuthService) Validate(ctx context.Context, accessToken string) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.cfg.JwtSecret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	uidFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid subject")
	}
	return uint64(uidFloat), nil
}

func (service *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return service.repository.RevokeRefreshToken(ctx, refreshToken)
}

func (service *AuthService) issueTokens(ctx context.Context, userID uint64) (pbmodel.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(service.cfg.AccessTokenTTLMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": accessExp.Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := access.SignedString([]byte(service.cfg.JwtSecret))
	if err != nil {
		return pbmodel.TokenPair{}, fmt.Errorf("could not sign access token: %w", err)
	}
	refreshRaw := fmt.Sprintf("%d:%d:%d", userID, now.UnixNano(), service.cfg.JwtSecret)
	h := sha256.Sum256([]byte(refreshRaw))
	refreshStr := hex.EncodeToString(h[:])
	refresh := model.RefreshToken{
		UserID:    userID,
		Token:     refreshStr,
		ExpiresAt: now.Add(time.Duration(service.cfg.RefreshTokenTTLDays) * 24 * time.Hour),
		CreatedAt: now,
	}
	if err := service.repository.SaveRefreshToken(ctx, refresh); err != nil {
		return pbmodel.TokenPair{}, err
	}

	return pbmodel.TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}
