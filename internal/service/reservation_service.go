package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type ReservationService struct {
	ReservationRepository domain.ReservationRepositoryInterface
}

func (repo *ReservationService) CreateReservationService(ctx context.Context) error {
	log := pkg.GetLogger()

	reservation, ok := utils.GetContextValue[domain.CreateReservation](
		ctx,
		utils.CreateReservationKey,
	)
	if !ok {
		log.Error("Error: Invalid reservation type in context")
		return pkg.NewBadRequestError("invalid reservation type in context")
	}

	log.Info("Processing reservation for desk: %s", reservation.DeskId)

	isReservationMade, err := checkReservationMade(ctx, repo, reservation)
	if err != nil {
		return err
	}

	if isReservationMade.Id == "" {
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
) (domain.Reservation, error) {
	log := pkg.GetLogger()

	reservation, err := repo.ReservationRepository.FindReservation(ctx, reservationData)
	if err != nil {
		log.Error("Error to find reservation: %v", err)
		return domain.Reservation{}, pkg.NewInternalServerError("failed to find reservation", err)
	}

	if reservation == nil {
		log.Info("No existing reservation found for desk: %s", reservationData.DeskId)
		return domain.Reservation{}, nil
	}

	return *reservation, nil
}
