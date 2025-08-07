package service

import (
	"context"
	"testing"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/utils"
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
        // Setup
        ctx := context.Background()
        user := domain.CreateUser{
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "password123",
        }
        ctx = utils.SetContextValue(ctx, utils.CreateUserKey, user)

        // mock := &mockDB{
        //     findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
        //         return nil, nil // User doesn't exist
        //     },
        //     saveUserFunc: func(ctx context.Context, user domain.CreateUser) error {
        //         return nil // Save successful
        //     },
        // }

        // Execute
        err := CreateUserService(ctx)

        // Verify
        if err != nil {
            t.Errorf("Expected no error, got %v", err)
        }
    })

    t.Run("should handle invalid context value", func(t *testing.T) {
        // Setup
        ctx := context.Background()

        // Execute
        err := CreateUserService(ctx)

        // Verify
        if err == nil {
            t.Error("Expected error for invalid context, got nil")
        }
    })

    t.Run("should handle existing user", func(t *testing.T) {
        // Setup
        ctx := context.Background()
        user := domain.CreateUser{
            Name:     "Test User",
            Email:    "existing@example.com",
            Password: "password123",
        }
        ctx = utils.SetContextValue(ctx, utils.CreateUserKey, user)

        mock := &mockDB{
            findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
                return &domain.User{Email: email}, nil // User exists
            },
        }

        // Execute
        err := CreateUserService(ctx)

        // Verify
        if err == nil {
            t.Error("Expected error for existing user, got nil")
        }
    })

    t.Run("should handle database error", func(t *testing.T) {
        // Setup
        ctx := context.Background()
        user := domain.CreateUser{
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "password123",
        }
        ctx = utils.SetContextValue(ctx, utils.CreateUserKey, user)

        mock := &mockDB{
            findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
                return nil, pkg.NewInternalServerError("database error", nil)
            },
        }

        // Execute
        err := CreateUserService(ctx)

        // Verify
        if err == nil {
            t.Error("Expected error for database failure, got nil")
        }
    })
}

func TestCheckExistingUser(t *testing.T) {
    t.Run("should return nil for non-existing user", func(t *testing.T) {
        // Setup
        ctx := context.Background()
        email := "nonexistent@example.com"
        mock := &mockDB{
            findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
                return nil, nil
            },
        }

        // Execute
        err := checkExistingUser(ctx, mock, email)

        // Verify
        if err != nil {
            t.Errorf("Expected no error, got %v", err)
        }
    })

    t.Run("should return error for existing user", func(t *testing.T) {
        // Setup
        ctx := context.Background()
        email := "existing@example.com"
        mock := &mockDB{
            findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
                return &domain.User{Email: email}, nil
            },
        }

        // Execute
        err := checkExistingUser(ctx, mock, email)

        // Verify
        if err == nil {
            t.Error("Expected error for existing user, got nil")
        }
    })

    t.Run("should handle database error", func(t *testing.T) {
        // Setup
        ctx := context.Background()
        email := "test@example.com"
        mock := &mockDB{
            findUserByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
                return nil, pkg.NewInternalServerError("database error", nil)
            },
        }

        // Execute
        err := checkExistingUser(ctx, mock, email)

        // Verify
        if err == nil {
            t.Error("Expected error for database failure, got nil")
        }
    })
}
