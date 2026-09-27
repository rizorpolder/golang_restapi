package model

import "time"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type User struct {
	ID         uint64    `json:"id"`
	Username   string    `json:"username"`
	Login      string    `json:"login"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	Age        uint32    `json:"age"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateUser struct {
	Login      string
	Email      string
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Age        uint32
}
