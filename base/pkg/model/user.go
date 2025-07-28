package model

import "base/api/auth"

type User struct {
	Id           int
	Login        string
	Password     string
	PasswordHash string
}

type AuthResponse struct {
	User  *User
	Token string
}

func NewUserFromGrpc(u *auth.User) *User {
	return &User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func NewAuthResponseFromGrpc(r *auth.AuthResponse) *AuthResponse {
	return &AuthResponse{
		User:  NewUserFromGrpc(r.User),
		Token: r.Token,
	}
}

func (u *User) Grpc() *auth.User {
	return &auth.User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func (r *AuthResponse) Grpc() *auth.AuthResponse {
	return &auth.AuthResponse{
		User:  r.User.Grpc(),
		Token: r.Token,
	}
}