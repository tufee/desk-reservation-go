package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/infra"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func CreateUserService(ctx context.Context) error {
	log := pkg.GetLogger()

	user, ok := utils.GetContextValue[domain.CreateUser](ctx, utils.CreateUserKey)
	if !ok {
		log.Error("Error: Invalid user type in context")
		return pkg.NewBadRequestError("invalid user type in context")
	}

	log.Info("Processing user creation for email: %s", user.Email)

	db, err := infra.InitializeDB()
	if err != nil {
		return err
	}

	if err := checkExistingUser(ctx, db, user.Email); err != nil {
		return err
	}

	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return pkg.NewInternalServerError("error processing user data", err)
	}

	user.Password = hashedPassword

	if err := db.SaveUser(ctx, user); err != nil {
		log.Error("Error saving user to database: %v", err)
		return err
	}

	log.Info("Successfully created user with email: %s", user.Email)
	return nil
}

func checkExistingUser(ctx context.Context, db *infra.Db, email string) error {
	log := pkg.GetLogger()

	existingUser, err := db.FindUserByEmail(ctx, email)
	if err != nil {
		log.Error("Error checking existing user: %v", err)
		return pkg.NewInternalServerError("failed to check existing user", err)
	}

	if existingUser != nil {
		log.Warn("User with email %s already exists", email)
		return pkg.NewBadRequestError("user already exists")
	}
	return nil
}
