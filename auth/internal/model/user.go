package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

var UserRoles = struct {
	Admin  UserRole
	Client UserRole
}{
	Admin:  "admin",
	Client: "client",
}

type User struct {
	Id           int
	Login        string
	Password     string `db:"-"`
	PasswordHash string
	Role         UserRole
}

type AuthResponse struct {
	Token string
}

type AuthClaims struct {
	UserId int
	Login  string
	Role   UserRole
	jwt.RegisteredClaims
}