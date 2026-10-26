package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"golang_restapi/internal/auth/config"
	"golang_restapi/internal/auth/model"
	"time"
)

type AuthService struct {
	repo   Repository
	cfg    *config.Config
	logger *zerolog.Logger
}

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, userId uint64) error
	GetUserByID(ctx context.Context, id uint64) (model.User, error)
	GetUserByLoginOrEmail(ctx context.Context, loginOrEmail string) (model.User, error)
	SaveRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

func New(repo Repository, cfg *config.Config, logger *zerolog.Logger) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, logger: logger}
}

func (s *AuthService) Register(ctx context.Context, req model.Register) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}

	user := model.User{
		Login:        req.Login,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req model.Login) (model.TokenPair, error) {
	user, err := s.repo.GetUserByLoginOrEmail(ctx, req.LoginOrEmail)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid credentials")
	}
	return s.issueTokens(ctx, user.ID)
}
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("invalid refresh token")
	}
	if rt.RevokedAt == nil {
		return model.TokenPair{}, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		return model.TokenPair{}, fmt.Errorf("refresh token is expired or revoked")
	}
	return s.issueTokens(ctx, rt.ID)
}
func (s *AuthService) Validate(ctx context.Context, accessToken string) (uint64, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
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

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint64) (model.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(time.Duration(s.cfg.AccessTokenTTLMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": accessExp.Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := access.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("could not sign access token: %w", err)
	}
	refreshRaw := fmt.Sprintf("%d:%d:%d", userID, now.UnixNano(), s.cfg.JWTSecret)
	h := sha256.Sum256([]byte(refreshRaw))
	refreshStr := hex.EncodeToString(h[:])
	refresh := model.RefreshToken{
		UserID:    userID,
		Token:     refreshStr,
		ExpiresAt: now.Add(time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.repo.SaveRefreshToken(ctx, refresh); err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{AccessToken: accessStr, RefreshToken: refreshStr}, nil
}
