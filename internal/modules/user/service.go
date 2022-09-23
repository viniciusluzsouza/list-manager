package user

import (
	"log"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
)

type (
	Service interface {
		Save(user *models.User) error
		Get(id uint64) (*models.User, error)
		Update(id uint64, user *models.User) (*models.User, error)
	}

	service struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s service) Save(user *models.User) error {
	err := s.checkIfLoginExists(user.Login)
	if err != nil {
		return err
	}

	err = user.HashPassword()
	if err != nil {
		log.Printf("Error hashing user password: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error hashing user password")
	}

	err = s.repository.Save(user)
	if err != nil {
		log.Printf("Error saving user: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error saving user")
	}

	user.Password = ""
	return nil
}

func (s service) Get(id uint64) (*models.User, error) {
	user, err := s.repository.Get(id)
	if err != nil {
		log.Printf("Error getting user: %s\n", err.Error())
		return nil, apperrors.NewInternalError("Internal error getting user")
	}

	if user == nil {
		return nil, apperrors.NewNotFoundError("user", id)
	}

	user.Password = ""
	return user, nil
}

func (s service) Update(id uint64, user *models.User) (*models.User, error) {
	dbUser, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// Only update email and name, for while
	dbUser.Email = user.Email
	dbUser.Name = user.Name

	err = s.repository.Update(dbUser)
	if err != nil {
		log.Printf("Error updating user: %s\n", err.Error())
		return nil, apperrors.NewInternalError("Internal error updating user")
	}

	dbUser.Password = ""
	return dbUser, nil
}

func (s service) checkIfLoginExists(login string) error {
	exists, err := s.repository.ExistsByLogin(login)
	if err != nil {
		return apperrors.NewInternalError("Internal error checking if login already registered")
	}

	if exists {
		return apperrors.NewLoginAlreadyRegisteredError(login)
	}

	return nil
}
