package middlewares

import (
	"errors"
	"net/http"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/internal/auth"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/constants"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			err := errors.New("missing Authorization token")
			context.JSON(http.StatusUnauthorized, models.NewHttpError(err))
			context.Abort()
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			context.JSON(http.StatusUnauthorized, models.NewHttpError(err))
			context.Abort()
			return
		}

		setClaimInContext(context, claims)
		context.Next()
	}
}

func PublicAuthenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token != "" {
			claims, err := auth.GetClaims(token)
			if err == nil {
				setClaimInContext(context, claims)
			}
		}
		context.Next()
	}
}

func setClaimInContext(context *gin.Context, claim *auth.JWTClaim) {
	context.Set(constants.CtxUserKey, claim.ID)
}
