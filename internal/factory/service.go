package factory

import (
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/auth"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/item"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/list"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
)

func NewAuthService(repository user.Repository) auth.Service {
	return auth.NewService(repository)
}

func NewUserService(repository user.Repository) user.Service {
	return user.NewService(repository)
}

func NewListService(repository list.Repository) list.Service {
	return list.NewService(repository)
}

func NewItemService(
	repository item.Repository,
	listRepository list.Repository,
	userRepository user.Repository,
) item.Service {
	return item.NewService(repository, listRepository, userRepository)
}
