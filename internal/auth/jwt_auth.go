package auth

import (
	"errors"
	"os"
	"time"

	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/constants"
	"git.vibbra.com.br/vinicius-1663626255/vibbra-list-manager/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(constants.TokenExpirationTime)
	claims := &JWTClaim{
		ID:    user.ID,
		Email: user.Email,
		Login: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateToken(signedToken string) (*JWTClaim, error) {
	claims, err := GetClaims(signedToken)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func GetClaims(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("cannot parse token claims")
	}

	return claims, nil
}

func InitJWTAuth() {
	key := os.Getenv("JWT_SECRET")
	jwtKey = []byte(key)
}

func ValidateTokenSSO(token string, login string) error {
	claims, err := ValidateToken(token)
	if err != nil {
		return err
	}

	if claims.Login != login {
		return errors.New("invalid login sso")
	}

	return nil
}
