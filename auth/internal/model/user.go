package model

type User struct {
	Id       int
	Login    string
	Password string
}

type AuthResponse struct {
	User  *User
	Token string
}