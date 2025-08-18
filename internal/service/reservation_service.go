package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type ReservationService struct {
	ReservationRepository domain.ReservationRepositoryInterface
}

func (repo *ReservationService) CreateReservationService(
	ctx context.Context,
	reservation domain.CreateReservation,
) error {
	log := pkg.GetLogger()

	log.Info("Processing reservation for desk: %s", reservation.DeskId)

	isReservationMade, err := checkReservationMade(ctx, repo, reservation)
	if err != nil {
		return err
	}

	if isReservationMade == nil {
		if err := repo.ReservationRepository.SaveReservation(ctx, reservation); err != nil {
			log.Error("Error saving user to database: %v", err)
			return err
		}
		log.Info("reservation created successfully")
		return nil
	}
	return pkg.NewBadRequestError("desk is unavailable")
}

func checkReservationMade(
	ctx context.Context,
	repo *ReservationService,
	reservationData domain.CreateReservation,
) (*domain.Reservation, error) {
	log := pkg.GetLogger()

	reservation, err := repo.ReservationRepository.FindReservation(ctx, reservationData)
	if err != nil {
		log.Error("Error to find reservation: %v", err)
		return nil, pkg.NewInternalServerError("failed to find reservation", err)
	}

	if reservation == nil {
		log.Info("No existing reservation found for desk: %s", reservationData.DeskId)
		return nil, nil
	}

	return reservation, nil
}
