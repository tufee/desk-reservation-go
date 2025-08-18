package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tufee/desk-reservation-go/internal/domain"
)

type reservationRepo struct {
	FindReservationFunc func(ctx context.Context, reservation domain.CreateReservation) (*domain.Reservation, error)
	SaveReservationFunc func(ctx context.Context, reservation domain.CreateReservation) error
}

func (r *reservationRepo) FindReservation(
	ctx context.Context,
	reservation domain.CreateReservation,
) (*domain.Reservation, error) {
	return r.FindReservationFunc(ctx, reservation)
}

func (r *reservationRepo) SaveReservation(
	ctx context.Context,
	reservation domain.CreateReservation,
) error {
	return r.SaveReservationFunc(ctx, reservation)
}

func TestCreateReservationService(t *testing.T) {
	t.Run("should create reservation successfully", func(t *testing.T) {
		parsedTime, _ := time.Parse(time.RFC3339, "2025-06-05T00:54:07Z")

		data := domain.CreateReservation{
			DeskId: "48b8c429-be55-470f-a245-651fc3c75a6b",
			UserId: "1a162e27-45ff-4632-817a-a79e88c8f878",
			Date:   parsedTime,
		}

		mock := &reservationRepo{
			FindReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) (*domain.Reservation, error) {
				return nil, nil
			},
			SaveReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) error {
				return nil
			},
		}

		ctx := context.Background()
		reservation := ReservationService{ReservationRepository: mock}
		err := reservation.CreateReservationService(ctx, data)

		assert.NoError(t, err, "should not return error")
		assert.Equal(t, nil, err, "should create reservation successfully")
	})

	t.Run("should return desk unavailable", func(t *testing.T) {
		parsedTime, _ := time.Parse(time.RFC3339, "2025-06-05T00:54:07Z")

		data := domain.CreateReservation{
			DeskId: "48b8c429-be55-470f-a245-651fc3c75a6b",
			UserId: "1a162e27-45ff-4632-817a-a79e88c8f878",
			Date:   parsedTime,
		}

		mock := &reservationRepo{
			FindReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) (*domain.Reservation, error) {
				return &domain.Reservation{}, nil
			},
			SaveReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) error {
				return nil
			},
		}

		ctx := context.Background()
		reservation := ReservationService{ReservationRepository: mock}
		err := reservation.CreateReservationService(ctx, data)

		assert.Error(t, err, "should return erro")
		assert.Equal(t, err.Error(), "desk is unavailable", "should return correct message")
	})

	t.Run("should fail to find reservation", func(t *testing.T) {
		parsedTime, _ := time.Parse(time.RFC3339, "2025-06-05T00:54:07Z")

		data := domain.CreateReservation{
			DeskId: "48b8c429-be55-470f-a245-651fc3c75a6b",
			UserId: "1a162e27-45ff-4632-817a-a79e88c8f878",
			Date:   parsedTime,
		}

		mock := &reservationRepo{
			FindReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) (*domain.Reservation, error) {
				return nil, errors.New("random error")
			},
			SaveReservationFunc: func(
				ctx context.Context,
				reservation domain.CreateReservation,
			) error {
				return nil
			},
		}

		ctx := context.Background()
		reservation := ReservationService{ReservationRepository: mock}
		err := reservation.CreateReservationService(ctx, data)

		assert.Error(t, err, "should return erro")
		assert.Equal(
			t,
			err.Error(),
			"failed to find reservation: random error",
			"should fail to find reservation",
		)
	})
}
