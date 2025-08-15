package service

import (
	"base/pkg/model"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJwtToken(token string) (*model.AuthClaims, error) {
	claims := &model.AuthClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token, claims)
	return claims, err
}