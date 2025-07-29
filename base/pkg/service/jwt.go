package service

import (
	"base/pkg/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authClaims struct {
	UserId int
	Login  string
	jwt.RegisteredClaims
}

func EncodeJwtToken(user *model.User, secretKey string) (string, error) {
	claims := authClaims{
		UserId: user.Id,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func DecodeJwtToken(token string, secretKey string) (*authClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&authClaims{},
		func(*jwt.Token) (any, error) {
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*authClaims)
	if !ok {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}