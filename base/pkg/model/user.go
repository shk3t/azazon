package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type userRole string

var UserRoles = struct {
	Admin  userRole
	Client userRole
}{
	Admin:  "admin",
	Client: "client",
}

type User struct {
	Id           int
	Login        string
	Password     string
	PasswordHash string
	Role         userRole
}

type AuthResponse struct {
	Token string
}

type AuthClaims struct {
	UserId int
	Login  string
	Role   userRole
	jwt.RegisteredClaims
}