package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/infra"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func CreateReservationService(ctx context.Context) error {
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

	db, err := infra.InitializeDB()
	if err != nil {
		return pkg.NewInternalServerError("failed to initialize database", err)
	}

	isReservationMade, err := checkReservationMade(ctx, db, reservation)
	if err != nil {
		return err
	}

	if isReservationMade.Id != "xx" {
		return pkg.NewBadRequestError("desk is unavailable")
	}

	if err := db.SaveReservation(ctx, reservation); err != nil {
		log.Error("Error saving user to database: %v", err)
		return pkg.NewBadRequestError("error to create reservation")
	}

	log.Info("reservation created successfully")
	return nil
}

func checkReservationMade(
	ctx context.Context,
	db *infra.Db,
	reservationData domain.CreateReservation,
) (domain.Reservation, error) {
	log := pkg.GetLogger()

	reservation, err := db.FindReservation(ctx, reservationData)
	if err != nil {
		log.Error("Error to find reservation: %v", err)
		return domain.Reservation{}, pkg.NewInternalServerError("failed to find reservation", err)
	}

	return *reservation, nil
}
