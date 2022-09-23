package list

import (
	"errors"
	"net/http"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/apperrors"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/constants"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/utils"
	"github.com/gin-gonic/gin"
)

type (
	Handler interface {
		Save(c *gin.Context)
		Get(c *gin.Context)
		Delete(c *gin.Context)
	}

	handler struct {
		service Service
	}
)

func NewHandler(service Service) Handler {
	return &handler{service}
}

func (h handler) Save(c *gin.Context) {
	listDTO, err := getListFromRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	userID := c.GetUint64(constants.CtxUserKey)

	list := models.NewListFromDTO(listDTO)
	err = h.service.Save(list, userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.NewHttpError(err))
		return
	}

	c.IndentedJSON(http.StatusCreated, models.NewListDTO(list))
}

func (h handler) Get(c *gin.Context) {
	id, err := utils.GetIDFromRequest(c, "list_id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	list, err := h.service.Get(id)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, models.NewListDTO(list))
}

func (h handler) Delete(c *gin.Context) {
	id, err := utils.GetIDFromRequest(c, "list_id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func getListFromRequest(c *gin.Context) (*models.ListDTO, error) {
	var list models.ListDTO

	if err := c.BindJSON(&list); err != nil {
		return nil, errors.New("invalid body request format")
	}

	return &list, nil
}
