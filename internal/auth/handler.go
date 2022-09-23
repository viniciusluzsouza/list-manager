package auth

import (
	"errors"
	"net/http"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/gin-gonic/gin"
)

type (
	Handler interface {
		Authenticate(c *gin.Context)
		AuthenticateSSO(c *gin.Context)
	}

	handler struct {
		service Service
	}
)

func NewHandler(service Service) Handler {
	return &handler{service}
}

func (h handler) Authenticate(c *gin.Context) {
	authRequest, err := getAuthBodyRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	authResponse, err := h.service.Authenticate(authRequest)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, authResponse)
}

func (h handler) AuthenticateSSO(c *gin.Context) {
	authRequest, err := getAuthSSOBodyRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	authResponse, err := h.service.AuthenticateSSO(authRequest)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, authResponse)
}

func (h handler) handleAuthError(c *gin.Context, err error) {
	switch err.(type) {
	case *apperrors.UserLoginError:
		c.IndentedJSON(http.StatusUnauthorized, models.NewHttpError(err))
	default:
		c.IndentedJSON(http.StatusInternalServerError, models.NewHttpError(err))
	}
}

func getAuthBodyRequest(c *gin.Context) (*models.AuthRequest, error) {
	var authRequest models.AuthRequest

	if err := c.BindJSON(&authRequest); err != nil {
		return nil, errors.New("invalid body request format")
	}

	return &authRequest, nil
}

func getAuthSSOBodyRequest(c *gin.Context) (*models.AuthRequestSSO, error) {
	var authRequest models.AuthRequestSSO

	if err := c.BindJSON(&authRequest); err != nil {
		return nil, errors.New("invalid body request format")
	}

	return &authRequest, nil
}
