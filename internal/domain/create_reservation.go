package domain

import (
	"time"
)

type CreateReservation struct {
	DeskId string    `json:"desk_id" validate:"required"`
	UserId string    `json:"user_id" validate:"required"`
	Date   time.Time `json:"date"    validate:"required"`
}
