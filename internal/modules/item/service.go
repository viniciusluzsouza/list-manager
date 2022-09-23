package item

import (
	"log"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/list"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
)

type (
	Service interface {
		Save(item *models.Item) error
		GetItemsFromList(listID uint64) (*[]models.Item, error)
		Update(user *models.Item) error
		Delete(listID uint64, itemID uint64) error
	}

	service struct {
		repository     Repository
		listRepository list.Repository
		userRepository user.Repository
	}
)

func NewService(
	repository Repository,
	listRepository list.Repository,
	userRepository user.Repository,
) Service {
	return &service{repository, listRepository, userRepository}
}

func (s service) Save(item *models.Item) error {
	err := s.checkIfListExists(item.ListID)
	if err != nil {
		return err
	}

	err = s.checkIfUserExists(*item.UserID)
	if err != nil {
		return err
	}

	err = s.repository.Save(item)
	if err != nil {
		log.Printf("Error saving item: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error saving item")
	}

	return nil
}

func (s service) GetItemsFromList(listID uint64) (*[]models.Item, error) {
	list, err := s.listRepository.Get(listID)
	if err != nil {
		log.Printf("Error getting list: %s\n", err.Error())
		return nil, apperrors.NewInternalError("Internal error getting items from list")
	}

	if list == nil {
		return nil, apperrors.NewNotFoundError("list", listID)
	}

	return s.repository.GetItemsFromList(list.ID)
}

func (s service) Update(item *models.Item) error {
	err := s.checkIfItemExistsInList(item.ListID, item.ID)
	if err != nil {
		return err
	}

	err = s.checkIfUserExists(*item.UserID)
	if err != nil {
		return err
	}

	err = s.repository.Update(item)
	if err != nil {
		log.Printf("Error updating item: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error updating item")
	}

	return nil
}

func (s service) Delete(listID uint64, itemID uint64) error {
	err := s.checkIfItemExistsInList(listID, itemID)
	if err != nil {
		return err
	}

	err = s.repository.Delete(itemID)
	if err != nil {
		log.Printf("Error deletting item: %s\n", err.Error())
		return apperrors.NewInternalError("Internal error deleting item")
	}

	return nil
}

func (s service) checkIfItemExistsInList(listID uint64, itemID uint64) error {
	err := s.checkIfListExists(listID)
	if err != nil {
		return err
	}

	isItemOnList, err := s.repository.IsItemInList(listID, itemID)
	if err != nil {
		return apperrors.NewInternalError("Internal error checking item in list")
	}

	if !isItemOnList {
		return apperrors.NewItemNotFoundInListError(itemID, listID)
	}

	return nil
}

func (s service) checkIfUserExists(userID uint64) error {
	exists, err := s.userRepository.Exists(userID)
	if err != nil {
		return apperrors.NewInternalError("Internal error checking if user exists")
	}

	if !exists {
		return apperrors.NewNotFoundError("user_id", userID)
	}

	return nil
}

func (s service) checkIfListExists(listID uint64) error {
	exists, err := s.listRepository.Exists(listID)
	if err != nil {
		return apperrors.NewInternalError("Internal error checking if list exists")
	}

	if !exists {
		return apperrors.NewNotFoundError("list", listID)
	}

	return nil
}
