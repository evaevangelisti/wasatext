package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func (service *UserService) GetUsers(q string, authenticatedUserID uuid.UUID) ([]models.User, error) {
	return service.Repository.GetUsers(q, authenticatedUserID)
}

func (service *UserService) DoLogin(username string) (*models.User, bool, error) {
	user, err := service.Repository.GetUserByUsername(username)
	if err != nil {
		return nil, false, err
	}

	if user != nil {
		return user, false, nil
	}

	userID, err := service.Repository.CreateUser(username)
	if err != nil {
		return nil, false, err
	}

	user, err = service.Repository.GetUserByID(userID)
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (service *UserService) UpdateUsername(userID uuid.UUID, username string) (*models.User, error) {
	existingUser, err := service.Repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if existingUser != nil && existingUser.ID != userID {
		return nil, errors.ErrConflict
	}

	err = service.Repository.UpdateUsername(userID, username)
	if err != nil {
		return nil, err
	}

	user, err := service.Repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateProfilePicture(userID uuid.UUID, profilePicture string) (*models.User, error) {
	err := service.Repository.UpdateProfilePicture(userID, profilePicture)
	if err != nil {
		return nil, err
	}

	user, err := service.Repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
