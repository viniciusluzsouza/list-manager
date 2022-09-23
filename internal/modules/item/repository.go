package item

import (
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Save(item *models.Item) error
		GetItemsFromList(listID uint64) (*[]models.Item, error)
		Update(item *models.Item) error
		Delete(id uint64) error
		IsItemInList(listID uint64, itemID uint64) (bool, error)
	}

	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r repository) Save(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r repository) GetItemsFromList(listID uint64) (*[]models.Item, error) {
	var items []models.Item
	err := r.db.Where(&models.Item{ListID: listID}).Find(&items).Error
	return &items, err
}

func (r repository) Update(item *models.Item) error {
	return r.db.Model(item).Updates(item).Error
}

func (r repository) Delete(id uint64) error {
	return r.db.Delete(&models.Item{}, id).Error
}

func (r repository) IsItemInList(listID uint64, itemID uint64) (bool, error) {
	var exists bool
	err := r.db.Model(&models.Item{}).
		Select("count(*) > 0").
		Where("id = ? and list_id = ?", itemID, listID).
		Find(&exists).
		Error
	return exists, err
}
