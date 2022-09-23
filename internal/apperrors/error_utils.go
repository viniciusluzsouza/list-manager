package apperrors

import (
	"net/http"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/gin-gonic/gin"
)

func HandleServiceError(c *gin.Context, err error) {
	switch err.(type) {
	case *NotFoundError:
		c.IndentedJSON(http.StatusNotFound, models.NewHttpError(err))
	case *ObjectInInvalidStateError:
		c.IndentedJSON(http.StatusBadRequest, models.NewHttpError(err))
	case *LoginAlreadyRegistered:
		c.IndentedJSON(http.StatusConflict, models.NewHttpError(err))
	default:
		c.IndentedJSON(http.StatusInternalServerError, models.NewHttpError(err))
	}
}
