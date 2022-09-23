package factory

import (
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/auth"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/item"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/list"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
)

func NewAuthHandler(service auth.Service) auth.Handler {
	return auth.NewHandler(service)
}

func NewUserHandler(service user.Service) user.Handler {
	return user.NewHandler(service)
}

func NewListHandler(service list.Service) list.Handler {
	return list.NewHandler(service)
}

func NewItemHandler(service item.Service) item.Handler {
	return item.NewHandler(service)
}
