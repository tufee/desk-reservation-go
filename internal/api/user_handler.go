package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/service"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var data domain.CreateUser

	if err := parseAndValidateRequest(r, &data, w); err != nil {
		return
	}

	user := buildUserFromRequest(data)
	ctx := context.WithValue(r.Context(), utils.CreateUserKey, user)

	if _, err := service.CreateUserService(ctx); err != nil {
		w.Header().Set("Content-Type", "application/json")

		switch e := err.(type) {
		case *pkg.BadRequestError:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"message": e.Error(),
			})

		case *pkg.InternalServerError:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"message": e.Error(),
			})

		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"message": "Internal server error",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User created successfully",
	})
}

func parseAndValidateRequest(r *http.Request, data *domain.CreateUser, w http.ResponseWriter) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		log.Println("CreateUser Validation error")
		respondWithValidationErrors(w, err.(validator.ValidationErrors))
		return err
	}
	return nil
}

func buildUserFromRequest(data domain.CreateUser) domain.User {
	now := time.Now()
	return domain.User{
		Id:         uuid.New().String(),
		Name:       data.Name,
		Email:      data.Email,
		Password:   data.Password,
		Created_at: now,
		Updated_at: now,
	}
}

func respondWithValidationErrors(w http.ResponseWriter, val validator.ValidationErrors) {
	errors := []map[string]string{}
	for _, err := range val {
		errors = append(errors, map[string]string{
			"field": err.Field(),
			"tag":   err.Tag(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errors,
	})
}
