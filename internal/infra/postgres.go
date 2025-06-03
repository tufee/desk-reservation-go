package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type Db struct {
	Conn *sqlx.DB
}

func InitializeDB() (*Db, error) {
	connection := os.Getenv("CONNECTION_STRING")

	if connection == "" {
		return nil, fmt.Errorf("missing DATABASE_URL in environment")
	}

	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &Db{Conn: db}, nil
}

func (db *Db) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
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

func (db *Db) SaveUser(ctx context.Context, user domain.CreateUser) error {
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

func (db *Db) FindReservation(
	ctx context.Context,
	reservation domain.CreateReservation,
) (*domain.Reservation, error) {
	query := `
        SELECT * FROM reservations
        WHERE desk_id = :desk_id
        AND DATE(date) = :date
        AND (status = 'pending' OR status = 'confirmed')
    `

	params := map[string]any{
		"desk_id": reservation.DeskId,
		"date": reservation.Date.Format(
			"2006-01-02",
		),
	}

	stmt, err := db.Conn.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, pkg.NewInternalServerError("failed to prepare statement", err)
	}
	defer stmt.Close()

	var result domain.Reservation

	err = stmt.GetContext(ctx, &result, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, pkg.NewInternalServerError("failed to find reservation", err)
	}

	return &result, nil
}

func (db *Db) SaveReservation(ctx context.Context, reservation domain.CreateReservation) error {
	query := `
	INSERT INTO reservations (desk_id, user_id, date)
	VALUES (:desk_id, :user_Id, :date)
	`
	_, err := db.Conn.NamedExecContext(ctx, query, reservation)
	if err != nil {
		return pkg.NewInternalServerError("failed to save reservation", err)
	}
	return nil
}
