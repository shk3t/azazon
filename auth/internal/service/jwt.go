package service

import (
	"auth/internal/config"
	"auth/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJwtToken(user model.User) (string, error) {
	claims := model.AuthClaims{
		UserId: user.Id,
		Login:  user.Login,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.AppName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.SecretKey))
}

func validateJwtToken(token string) error {
	_, err := jwt.ParseWithClaims(
		token,
		&model.AuthClaims{},
		func(*jwt.Token) (any, error) {
			return []byte(config.Env.SecretKey), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(config.AppName),
	)
	return err
}