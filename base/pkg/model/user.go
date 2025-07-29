package model

import "base/api/auth"

type User struct {
	Id           int
	Login        string
	Password     string
	PasswordHash string
}

type AuthResponse struct {
	Token string
}

func UserFromRegisterRequest(u *auth.RegisterRequest) *User {
	return &User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func UserFromLoginRequest(u *auth.LoginRequest) *User {
	return &User{
		Login:    u.Login,
		Password: u.Password,
	}
}

func (r *AuthResponse) RegisterResponse() *auth.RegisterResponse {
	return &auth.RegisterResponse{Token: r.Token}
}

func (r *AuthResponse) LoginResponse() *auth.LoginResponse {
	return &auth.LoginResponse{Token: r.Token}
}
