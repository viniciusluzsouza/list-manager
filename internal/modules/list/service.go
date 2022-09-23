package list

import (
	"log"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
)

type (
	Service interface {
		Save(list *models.List, userID uint64) error
		Get(id uint64) (*models.List, error)
		Delete(id uint64) error
	}

	service struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s service) Save(list *models.List, userID uint64) error {
	if userID != uint64(0) {
		list.Owner = &userID
	}

	err := s.repository.Save(list)
	if err != nil {
		log.Printf("Error saving list: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error saving list")
	}

	return nil
}

func (s service) Get(id uint64) (*models.List, error) {
	list, err := s.repository.Get(id)
	if err != nil {
		log.Printf("Error getting list: %s\n", err.Error())
		return nil, apperrors.NewInternalError("Internal error getting list")
	}

	if list == nil {
		return nil, apperrors.NewNotFoundError("list", id)
	}

	return list, nil
}

func (s service) Delete(id uint64) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}

	err = s.checkIfListIsEmpty(id)
	if err != nil {
		return err
	}

	err = s.repository.Delete(id)
	if err != nil {
		log.Printf("Error deletting list: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error deletting list")
	}

	return nil
}

func (s service) checkIfListIsEmpty(listID uint64) error {
	itemsCount, err := s.repository.CountItemsOnList(listID)
	if err != nil {
		return apperrors.NewInternalError("Internal error checking if list is empty")
	}

	if itemsCount > 0 {
		return apperrors.NewObjectInInvalidStateError("list is not empty")
	}

	return nil
}
