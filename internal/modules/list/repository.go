package list

import (
	"errors"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Save(list *models.List) error
		Get(id uint64) (*models.List, error)
		Delete(id uint64) error
		CountItemsOnList(id uint64) (int64, error)
		Exists(id uint64) (bool, error)
	}

	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r repository) Save(list *models.List) error {
	return r.db.Create(list).Error
}

func (r repository) Get(id uint64) (*models.List, error) {
	var list models.List
	err := r.db.First(&list, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &list, err
}

func (r repository) Delete(id uint64) error {
	return r.db.Delete(&models.List{}, id).Error
}

func (r repository) CountItemsOnList(id uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Item{}).Where("list_id = ?", id).Count(&count).Error
	return count, err
}

func (r repository) Exists(id uint64) (bool, error) {
	var exists bool
	err := r.db.Model(&models.List{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error
	return exists, err
}
