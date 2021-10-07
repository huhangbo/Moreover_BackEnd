package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("moreover")

var tokenExpireDuration = time.Hour * 24 * 7

type Claims struct {
	StuId string
	jwt.StandardClaims
}

func GenerateToken(stuId string) string {
	newClaims := Claims{
		stuId, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    "flying",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("tokenSign fail, err: %v\n", err)
	}
	return tokenString
}

func ParseToken(token string) *Claims {
	var newClaims = new(Claims)
	tmpToken, err := jwt.ParseWithClaims(token, newClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil
	}
	if tmpToken != nil {
		if tmpClaims, ok := tmpToken.Claims.(*Claims); ok && tmpToken.Valid {
			return tmpClaims
		}
	}
	return nil
}
