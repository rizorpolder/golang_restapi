package model

import "time"

type User struct {
	ID           uint64
	Login        string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        uint64
	UserID    uint64
	Token     string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
