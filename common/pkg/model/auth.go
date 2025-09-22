package model

import "github.com/golang-jwt/jwt/v5"

type AuthClaims struct {
	UserId int
	Login  string
	Role   string
	jwt.RegisteredClaims
}

func (c *AuthClaims) IsAdmin() bool {
	return c.Role == "admin"
}