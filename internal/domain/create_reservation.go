package domain

import (
	"time"
)

type CreateReservation struct {
	DeskId string    `json:"desk_id" db:"desk_id" validate:"required"`
	UserId string    `json:"user_id" db:"user_id" validate:"required"`
	Date   time.Time `json:"date" db:"date" validate:"required"`
}
