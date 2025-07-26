package domain

import (
	"context"
	"time"
)

type ReservationRepositoryInterface interface {
	FindReservation(ctx context.Context, reservation CreateReservation) (*Reservation, error)
	SaveReservation(ctx context.Context, reservation CreateReservation) error
}

type Reservation struct {
	Id        string    `json:"id"         db:"id"`
	DeskId    string    `json:"desk_id"    db:"desk_id"`
	UserId    string    `json:"user_id"    db:"user_id"`
	Date      time.Time `json:"date"       db:"date"`
	Status    string    `json:"status"     db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
