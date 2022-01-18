package service

import (
	"UserMicro/configs"
	"UserMicro/proto"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var key = []byte(configs.TOKEN_SECRET)

type CustomClaims struct {
	User *proto.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *proto.User) (string, error)
}

type TokenService struct {
}

func (tkn *TokenService) Decode(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return &CustomClaims{}, err
}

func (tkn *TokenService) Encode(user *proto.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()

	claims := CustomClaims{
		User:           user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
