package auth

import (
	"log"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/modules/user"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
)

type (
	Service interface {
		Authenticate(authRequest *models.AuthRequest) (*models.AuthResponse, error)
		AuthenticateSSO(authRequest *models.AuthRequestSSO) (*models.AuthResponse, error)
	}

	service struct {
		repository user.Repository
	}
)

func NewService(repository user.Repository) Service {
	return &service{repository}
}

func (s service) Authenticate(authRequest *models.AuthRequest) (*models.AuthResponse, error) {
	return s.authenticate(authRequest.Login, authRequest.Password, false)
}

func (s service) AuthenticateSSO(authRequest *models.AuthRequestSSO) (*models.AuthResponse, error) {
	err := ValidateTokenSSO(authRequest.APPToken, authRequest.Login)
	if err != nil {
		return nil, apperrors.NewUserSSOLoginError()
	}

	return s.authenticate(authRequest.Login, "", true)
}

func (s service) authenticate(login string, password string, isSSO bool) (*models.AuthResponse, error) {
	user, err := s.repository.GetByLogin(login)
	if err != nil {
		log.Printf("Error getting user to login: %s\n", err.Error())
		return nil, apperrors.NewInternalError("Internal error getting user to login")
	}

	if user == nil {
		return nil, apperrors.NewUserLoginError()
	}

	if !isSSO {
		err = user.CheckPassword(password)
		if err != nil {
			return nil, apperrors.NewUserLoginError()
		}
	}

	token, err := GenerateJWT(user)
	if err != nil {
		log.Printf("Error generating token for user %s: %s\n", user.Login, err.Error())
		return nil, apperrors.NewInternalError("Internal error generating token")
	}

	user.Password = ""
	response := models.AuthResponse{
		User:  user,
		Token: token,
	}

	return &response, nil
}
