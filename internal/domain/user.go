package domain

import (
	"context"
	"time"
)

type UserRepositoryInterface interface {
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	SaveUser(ctx context.Context, user CreateUser) error
}

type User struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
