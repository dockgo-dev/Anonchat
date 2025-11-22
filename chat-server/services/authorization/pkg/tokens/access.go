package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gox7/notify/services/authorization/models"
)

type (
	Claims struct {
		UserId int64  `json:"user_id"`
		Login  string `json:"login"`
		Email  string `json:"email"`
		jwt.RegisteredClaims
	}
)

func GenerateAccess(config *models.LocalConfig, userId int64, login string, email string) (string, error) {
	claims := Claims{
		UserId: userId, Login: login, Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Server.Password))
}

func CheckAccess(config *models.LocalConfig, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Password), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
