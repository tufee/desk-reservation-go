package infra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type UserRepositoryDb struct {
	Conn *sqlx.DB
}

func (db *UserRepositoryDb) FindUserByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`

	err := db.Conn.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, pkg.NewInternalServerError("failed to query user by email", err)
	}

	return &user, nil
}

func (db *UserRepositoryDb) SaveUser(ctx context.Context, user domain.CreateUser) error {
	query := `
	INSERT INTO users (name, email, password)
	VALUES (:name, :email, :password)
	ON CONFLICT (email) DO NOTHING
	`

	_, err := db.Conn.NamedExecContext(ctx, query, user)
	if err != nil {
		return pkg.NewInternalServerError("failed to save user", err)
	}

	return nil
}
