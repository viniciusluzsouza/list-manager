package user

import (
	"errors"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Save(user *models.User) error
		Get(id uint64) (*models.User, error)
		GetByLogin(login string) (*models.User, error)
		Update(user *models.User) error
		Exists(id uint64) (bool, error)
		ExistsByLogin(login string) (bool, error)
	}

	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r repository) Save(user *models.User) error {
	return r.db.Create(user).Error
}

func (r repository) Get(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r repository) GetByLogin(login string) (*models.User, error) {
	var user models.User
	err := r.db.Where(&models.User{Login: login}).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r repository) Update(user *models.User) error {
	return r.db.Model(user).Updates(user).Error
}

func (r repository) Exists(id uint64) (bool, error) {
	var exists bool
	err := r.db.Model(&models.User{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).
		Error
	return exists, err
}

func (r repository) ExistsByLogin(login string) (bool, error) {
	var exists bool
	err := r.db.Model(&models.User{}).
		Select("count(*) > 0").
		Where("login = ?", login).
		Find(&exists).
		Error
	return exists, err
}
