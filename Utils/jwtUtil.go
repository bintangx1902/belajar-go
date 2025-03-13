package Utils

import (
	"TestApp/Middlewares"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("190203")

type UserClaims struct {
	UserID           uint
	RegisteredClaims jwt.RegisteredClaims
}

func GenerateJWT(userID uint) (string, error) {
	claims := Middlewares.UserClaims{
		UserID: uint64(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("190203"))

}
