package service

import (
	"context"

	"github.com/tufee/desk-reservation-go/internal/domain"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

type UserService struct {
	UserRepository domain.UserRepositoryInterface
}

func (repo *UserService) CreateUserService(ctx context.Context, user domain.CreateUser) error {
	log := pkg.GetLogger()

	log.Info("Processing user creation for email: %s", user.Email)

	if err := checkExistingUser(ctx, repo, user.Email); err != nil {
		return err
	}

	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return pkg.NewInternalServerError("error processing user data", err)
	}

	user.Password = hashedPassword

	if err := repo.UserRepository.SaveUser(ctx, user); err != nil {
		log.Error("Error saving user to database: %v", err)
		return err
	}

	log.Info("Successfully created user with email: %s", user.Email)
	return nil
}

func checkExistingUser(ctx context.Context, repo *UserService, email string) error {
	log := pkg.GetLogger()

	existingUser, err := repo.UserRepository.FindUserByEmail(ctx, email)
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
