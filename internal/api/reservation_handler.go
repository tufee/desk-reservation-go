package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/infra"
	repo "github.com/tufee/desk-reservation-go/internal/infra/repository"
	"github.com/tufee/desk-reservation-go/internal/service"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	var data domain.CreateReservation

	if err := pkg.ParseAndValidateRequest(r, &data, w); err != nil {
		return
	}

	reservation := buildReservationFromRequest(data)
	ctx := context.Background()

	db, err := infra.InitializeDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	reservationRepository := &repo.ReservationRepositoryDb{Conn: db.Conn}
	reservationService := service.ReservationService{ReservationRepository: reservationRepository}

	if err := reservationService.CreateReservationService(ctx, reservation); err != nil {
		pkg.HandleHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Reservation created successfully",
	})
}

func buildReservationFromRequest(data domain.CreateReservation) domain.CreateReservation {
	return domain.CreateReservation{
		DeskId: data.DeskId,
		UserId: data.UserId,
		Date:   data.Date,
	}
}
