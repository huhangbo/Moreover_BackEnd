package jwt

import (
	"Moreover/internal/pkg/redis"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("moreover")

var tokenExpireDuration = time.Hour * 24 * 7

type Claims struct {
	StudentID string
	jwt.StandardClaims
}

func GenerateToken(studentID string) string {
	newClaims := Claims{
		studentID, jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    "flying",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("tokenSign fail, err: %v\n", err)
	}
	redis.DB.Set("token"+studentID, tokenString, tokenExpireDuration)
	return tokenString
}

func ParseToken(token string) *Claims {
	var newClaims = new(Claims)
	_, err := jwt.ParseWithClaims(token, newClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Printf("parseToken fail, err: %v\n", err)
		return nil
	}
	return newClaims
}
