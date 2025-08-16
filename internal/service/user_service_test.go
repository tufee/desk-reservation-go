package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type mockDB struct {
	findUserByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	saveUserFunc        func(ctx context.Context, user domain.CreateUser) error
}

func (m *mockDB) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.findUserByEmailFunc(ctx, email)
}

func (m *mockDB) SaveUser(ctx context.Context, user domain.CreateUser) error {
	return m.saveUserFunc(ctx, user)
}

func TestCreateUserService(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		ctx := context.Background()
		user := domain.CreateUser{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mock := &mockDB{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, nil
			},
			saveUserFunc: func(ctx context.Context, user domain.CreateUser) error {
				return nil
			},
		}

		userService := UserService{UserRepository: mock}
		err := userService.CreateUserService(ctx, user)

		assert.Equal(t, err, nil, "Error should be nil")
	})

	t.Run("should handle existing user", func(t *testing.T) {
		user := domain.CreateUser{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "password123",
		}

		mock := &mockDB{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{Email: email}, nil
			},
		}

		ctx := context.Background()

		userService := UserService{UserRepository: mock}
		err := userService.CreateUserService(ctx, user)

		assert.Error(t, err, "should return an error")
		assert.Equal(t, "user already exists", err.Error(), "should return correct error message")
	})

	t.Run("should handle database error", func(t *testing.T) {
		user := domain.CreateUser{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mock := &mockDB{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, pkg.NewInternalServerError(
					"failed to query user by email",
					errors.New("Error"),
				)
			},
		}

		ctx := context.Background()
		userService := UserService{UserRepository: mock}
		err := userService.CreateUserService(ctx, user)

		assert.Error(t, err, "should return an error")
		assert.Equal(
			t,
			"failed to query user by email: Error",
			err.Error(),
			"should return correct error message",
		)
	})
}

func TestCheckExistingUser(t *testing.T) {
	t.Run("should return nil for non-existing user", func(t *testing.T) {
		email := "nonexistent@example.com"
		mock := &mockDB{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, nil
			},
		}

		ctx := context.Background()
		userService := UserService{UserRepository: mock}
		user := checkExistingUser(ctx, &userService, email)

		assert.Equal(t, nil, user, "should return nil for non-existing user")
	})

	t.Run("should return error for existing user", func(t *testing.T) {
		email := "existing@example.com"
		mock := &mockDB{
			findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{Email: email}, nil
			},
		}

		ctx := context.Background()
		userService := UserService{UserRepository: mock}
		user := checkExistingUser(ctx, &userService, email)

		assert.Equal(t, "user already exists", user.Error(), "should return correct error message")
	})
}
