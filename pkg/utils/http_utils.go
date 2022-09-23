package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIDFromRequest(c *gin.Context, key string) (uint64, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		err = errors.New("invalid id")
	}

	return id, err
}
