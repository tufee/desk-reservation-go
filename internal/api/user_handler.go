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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var data domain.CreateUser

	if err := pkg.ParseAndValidateRequest(r, &data, w); err != nil {
		return
	}

	user := buildUserFromRequest(data)
	ctx := context.Background()

	db, err := infra.InitializeDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userRepository := &repo.UserRepositoryDb{Conn: db.Conn}
	userService := service.UserService{UserRepository: userRepository}

	if err := userService.CreateUserService(ctx, user); err != nil {
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
