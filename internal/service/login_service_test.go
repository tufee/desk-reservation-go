package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tufee/desk-reservation-go/internal/domain"
)

func TestLoginService(t *testing.T) {
	t.Run("should login successfully", func(t *testing.T) {
		mock := &userRepo{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{
					Email:    "test@test.com",
					Password: "$2a$10$gG9c609HyTVVIX09MLuTAOpiXpLxJhYRFLS8lUsgxcivhHVu2Uk5.",
				}, nil
			},
		}
		credentials := domain.Credentials{
			Email:    "test@test.com",
			Password: "senha",
		}

		context := context.Background()

		userService := LoginService{UserRepository: mock}
		token, _ := userService.LoginService(context, credentials)

		assert.IsType(t, "", token.Token, "should be a string")
	})

	t.Run("should return error to find user by email", func(t *testing.T) {
		mock := &userRepo{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, nil
			},
		}
		credentials := domain.Credentials{
			Email:    "test@test.com",
			Password: "senha",
		}

		context := context.Background()

		userService := LoginService{UserRepository: mock}
		_, err := userService.LoginService(context, credentials)

		assert.Equal(t, "user not found", err.Error(), "should return user not found")
	})

	t.Run("should return error to invalid password", func(t *testing.T) {
		mock := &userRepo{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{
					Email:    "test@test.com",
					Password: "senha",
				}, nil
			},
		}
		credentials := domain.Credentials{
			Email:    "test@test.com",
			Password: "senha",
		}

		context := context.Background()

		userService := LoginService{UserRepository: mock}
		_, err := userService.LoginService(context, credentials)

		assert.Equal(t, "invalid password", err.Error(), "should return invalid password")
	})
}
