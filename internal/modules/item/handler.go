package item

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
		GetByList(c *gin.Context)
		Update(c *gin.Context)
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
	listID, err := utils.GetIDFromRequest(c, "list_id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	itemDTO, err := getItemFromRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	item := models.NewItemFromDTO(itemDTO)
	item.ListID = listID
	err = h.service.Save(item)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, models.NewItemDTO(item))
}

func (h handler) GetByList(c *gin.Context) {
	listID, err := utils.GetIDFromRequest(c, "list_id")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	items, err := h.service.GetItemsFromList(listID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, models.NewItemsDTO(items))
}

func (h handler) Update(c *gin.Context) {
	listID, itemID, httpErr := getListIDAndItemIDFromRequest(c)
	if httpErr != nil {
		c.IndentedJSON(http.StatusBadRequest, httpErr)
		return
	}

	itemDTO, err := getItemFromRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
		return
	}

	item := models.NewItemFromDTO(itemDTO)
	item.ID = itemID
	item.ListID = listID

	err = h.service.Update(item)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, models.NewItemDTO(item))
}

func (h handler) Delete(c *gin.Context) {
	listID, itemID, httpErr := getListIDAndItemIDFromRequest(c)
	if httpErr != nil {
		c.IndentedJSON(http.StatusBadRequest, httpErr)
		return
	}

	err := h.service.Delete(listID, itemID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func getItemFromRequest(c *gin.Context) (*models.ItemDTO, error) {
	var item models.ItemDTO

	if err := c.BindJSON(&item); err != nil {
		return nil, errors.New("invalid body request format")
	}

	return &item, validateItemDTO(&item)
}

func getListIDAndItemIDFromRequest(c *gin.Context) (uint64, uint64, *models.HttpError) {
	var listID, itemID uint64

	listID, err := utils.GetIDFromRequest(c, "list_id")
	if err != nil {
		httpErr := models.NewHttpError(err)
		return listID, itemID, &httpErr
	}

	itemID, err = utils.GetIDFromRequest(c, "item_id")
	if err != nil {
		httpErr := models.NewHttpError(err)
		return listID, itemID, &httpErr
	}

	return listID, itemID, nil
}

func validateItemDTO(item *models.ItemDTO) error {
	if item.Title == "" {
		return errors.New("title cannot be empty")
	}

	if item.UserID == nil {
		return errors.New("user id cannot be null")
	}

	return nil
}
