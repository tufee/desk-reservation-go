package infra

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/tufee/desk-reservation-go/internal/domain"
)

func setupReservationRepositoryTestDB(t *testing.T) (*ReservationRepositoryDb, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := &ReservationRepositoryDb{Conn: sqlxDB}

	return db, mock
}

func TestFindReservation(t *testing.T) {
	db, mock := setupReservationRepositoryTestDB(t)
	ctx := context.Background()
	reservation := domain.CreateReservation{
		DeskId: "123",
		UserId: "456",
		Date:   time.Now(),
	}

	t.Run("should find reservation successfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "desk_id", "user_id", "date", "status"}).
			AddRow("789", reservation.DeskId, reservation.UserId, reservation.Date, "pending")

		mock.ExpectPrepare("SELECT (.+) FROM reservations").
			ExpectQuery().
			WithArgs(reservation.DeskId, reservation.Date.Format("2006-01-02")).
			WillReturnRows(rows)

		result, err := db.FindReservation(ctx, reservation)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Error("expected result to not be nil")
		}
		if result.DeskId != reservation.DeskId {
			t.Errorf("expected desk_id %s, got %s", reservation.DeskId, result.DeskId)
		}
	})

	t.Run("should return nil when reservation not found", func(t *testing.T) {
		mock.ExpectPrepare("SELECT (.+) FROM reservations").
			ExpectQuery().
			WithArgs(reservation.DeskId, reservation.Date.Format("2006-01-02")).
			WillReturnError(sql.ErrNoRows)

		result, err := db.FindReservation(ctx, reservation)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != nil {
			t.Error("expected result to be nil")
		}
	})
}

func TestSaveReservation(t *testing.T) {
	db, mock := setupReservationRepositoryTestDB(t)
	ctx := context.Background()
	reservation := domain.CreateReservation{
		DeskId: "123",
		UserId: "456",
		Date:   time.Now(),
	}

	t.Run("should save reservation successfully", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO reservations").
			WithArgs(reservation.DeskId, reservation.UserId, reservation.Date).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := db.SaveReservation(ctx, reservation)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("should handle db error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO reservations").
			WithArgs(reservation.DeskId, reservation.UserId, reservation.Date).
			WillReturnError(fmt.Errorf("db error"))

		err := db.SaveReservation(ctx, reservation)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

