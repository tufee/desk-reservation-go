package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{Message: message}
}

type InternalServerError struct {
	Message string
	Err     error
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *InternalServerError) Unwrap() error {
	return e.Err
}

func NewInternalServerError(message string, err error) *InternalServerError {
	return &InternalServerError{
		Message: message,
		Err:     err,
	}
}

func ParseAndValidateRequest[T any](
	r *http.Request,
	data *T,
	w http.ResponseWriter,
) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		log.Println("Validation struct error")
		RespondWithValidationErrors(w, err.(validator.ValidationErrors))
		return err
	}
	return nil
}

func RespondWithValidationErrors(w http.ResponseWriter, val validator.ValidationErrors) {
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

func HandleHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch e := err.(type) {
	case *BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"message": e.Error(),
		})

	case *InternalServerError:
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
}
