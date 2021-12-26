package model

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strings"
)

type CustomClaim struct {
	jwt.StandardClaims
	ID int `json:"id"`
}

func GenToken(id int, lifeTime int64, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: lifeTime,
		},
	})
	return token.SignedString([]byte(key))
}

func GetClaim(tokenString, key string) (*CustomClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(*CustomClaim)
	if !ok {
		return nil, errors.New("failed conviction")
	}
	return claim, nil
}

func ParseBearer(value string) string {
	value = strings.TrimSpace(value)
	isBearer := strings.HasPrefix("Bearer", value)
	words := strings.Split(value, " ")
	if len(words) != 2 || isBearer {
		return ""
	}
	return strings.TrimSpace(words[1])
}
