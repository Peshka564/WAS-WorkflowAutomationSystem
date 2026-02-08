package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId int64, JWTSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString(JWTSecret)
}
