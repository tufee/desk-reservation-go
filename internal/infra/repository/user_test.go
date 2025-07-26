package infra

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/tufee/desk-reservation-go/internal/domain"
)

func setupUserRepositoryTestDB(t *testing.T) (*UserRepositoryDb, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := &UserRepositoryDb{Conn: sqlxDB}

	return db, mock
}

func TestFindUserByEmail(t *testing.T) {
	db, mock := setupUserRepositoryTestDB(t)
	ctx := context.Background()
	email := "test@example.com"

	t.Run("should find user successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).
			AddRow("123", "Test User", email, "hashedpassword")

		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(email).
			WillReturnRows(rows)

		user, err := db.FindUserByEmail(ctx, email)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user to not be nil")
		}
		if user.Email != email {
			t.Errorf("expected email %s, got %s", email, user.Email)
		}
	})

	t.Run("should return nil when user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := db.FindUserByEmail(ctx, email)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user != nil {
			t.Error("expected user to be nil")
		}
	})
}

func TestSaveUser(t *testing.T) {
	db, mock := setupUserRepositoryTestDB(t)
	ctx := context.Background()
	user := domain.CreateUser{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	t.Run("should save user successfully", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := db.SaveUser(ctx, user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should handle db error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.Name, user.Email, user.Password).
			WillReturnError(fmt.Errorf("db error"))

		err := db.SaveUser(ctx, user)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

