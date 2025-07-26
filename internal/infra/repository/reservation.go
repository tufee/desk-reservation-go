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

type ReservationRepositoryDb struct {
	Conn *sqlx.DB
}

func (db *ReservationRepositoryDb) FindReservation(
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

func (db *ReservationRepositoryDb) SaveReservation(
	ctx context.Context,
	reservation domain.CreateReservation,
) error {
	query := `
	INSERT INTO reservations (desk_id, user_id, date)
	VALUES (:desk_id, :user_id, :date)
	`
	_, err := db.Conn.NamedExecContext(ctx, query, reservation)
	if err != nil {
		return pkg.NewInternalServerError("failed to save reservation", err)
	}
	return nil
}
