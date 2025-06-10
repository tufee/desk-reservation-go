package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/infra"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func LoginService(ctx context.Context) (*domain.LoginResponse, error) {
	log := pkg.GetLogger()

	credentials, ok := utils.GetContextValue[domain.Credentials](
		ctx,
		utils.LoginKey,
	)

	if !ok {
		return nil, pkg.NewBadRequestError("invalid credentials type in context")
	}

	log.Info("Processing login for: %s", credentials.Email)

	db, err := infra.InitializeDB()
	if err != nil {
		return nil, err
	}

	user, err := GetUserByEmail(ctx, db, credentials.Email)
	if err != nil {
		log.Error("error to find user by email: %v", err)
		return nil, err
	}

	isPasswordMatch := pkg.CheckPasswordHash(credentials.Password, user.Password)

	if !isPasswordMatch {
		log.Info("Invalid password for user: %s", credentials.Email)
		return nil, pkg.NewBadRequestError("invalid email or password")
	}

	token, err := pkg.GenerateJWT(user.Id, user.Email)
	if err != nil {
		log.Error("Error generating JWT token: %v", err)
		return nil, pkg.NewInternalServerError("failed to generate JWT token", err)
	}
	return &domain.LoginResponse{Token: token}, nil
}

func GetUserByEmail(ctx context.Context, db *infra.Db, email string) (*domain.User, error) {
	log := pkg.GetLogger()

	user, err := db.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.Info("User not found with email: %s", email)
		return nil, pkg.NewBadRequestError("user not found")
	}

	return user, nil
}
