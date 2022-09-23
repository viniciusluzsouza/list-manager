package factory

import (
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/item"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/list"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) user.Repository {
	return user.NewRepository(db)
}

func NewListRepository(db *gorm.DB) list.Repository {
	return list.NewRepository(db)
}

func NewItemRepository(db *gorm.DB) item.Repository {
	return item.NewRepository(db)
}
