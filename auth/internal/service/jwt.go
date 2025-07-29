package service

import (
	"auth/internal/config"
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

func encodeJwtToken(user *model.User) (string, error) {
	claims := authClaims{
		UserId: user.Id,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.AppName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.SecretKey))
}

func DecodeJwtToken(token string) (*authClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token,
		&authClaims{},
		func(*jwt.Token) (any, error) {
			return []byte(config.Env.SecretKey), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(config.AppName),
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
