package api

import (
	"encoding/json"
	"net/http"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/service"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var data domain.CreateUser

	if err := pkg.ParseAndValidateRequest(r, &data, w); err != nil {
		return
	}

	user := buildUserFromRequest(data)
	ctx := utils.SetContextValue(r.Context(), utils.CreateUserKey, user)

	if err := service.CreateUserService(ctx); err != nil {
		pkg.HandleHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User created successfully",
	})
}

func buildUserFromRequest(data domain.CreateUser) domain.CreateUser {
	return domain.CreateUser{
		Name:                 data.Name,
		Email:                data.Email,
		Password:             data.Password,
		PasswordConfirmation: data.PasswordConfirmation,
	}
}
