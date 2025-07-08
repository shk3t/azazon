package model

import api "base/api/go"

type User struct {
	Id       int
	Login    string
	Password string
}

type AuthResponse struct {
	User  *User
	Token string
}

func NewUserFromGrpc(u *api.User) *User {
	return &User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func NewAuthResponseFromGrpc(r *api.AuthResponse) *AuthResponse {
	return &AuthResponse{
		User:  NewUserFromGrpc(r.User),
		Token: r.Token,
	}
}

func (u *User) Grpc() *api.User {
	return &api.User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func (r *AuthResponse) Grpc() *api.AuthResponse {
	return &api.AuthResponse{
		User:  r.User.Grpc(),
		Token: r.Token,
	}
}