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

type Login struct {
	UserID       uint64
	LoginOrEmail string
	Password     string
}

type Register struct {
	Login    string
	Email    string
	Password string
}

type TokenPair struct {
	UserID       uint64
	AccessToken  string
	RefreshToken string
}
