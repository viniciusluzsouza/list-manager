//go:generate mockery --all --case=snake --output=./mocks

package user

import (
	"errors"
	"net/http"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/utils"
	"github.com/gin-gonic/gin"
)

type (
	Handler interface {
		Save(c *gin.Context)
		Get(c *gin.Context)
		Update(c *gin.Context)
	}

	handler struct {
		service Service
	}
)

func NewHandler(service Service) Handler {
	return &handler{service}
}

func (h handler) Save(c *gin.Context) {
	user, err := getUserFromRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	err = h.service.Save(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.NewHttpError(err))
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}

func (h handler) Get(c *gin.Context) {
	id, err := utils.GetIDFromRequest(c, "id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	user, err := h.service.Get(id)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (h handler) Update(c *gin.Context) {
	id, err := utils.GetIDFromRequest(c, "id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	user, err := getUserFromRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	updatedUser, err := h.service.Update(id, user)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func getUserFromRequest(c *gin.Context) (*models.User, error) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		return nil, errors.New("invalid body request format")
	}

	return &user, nil
}
