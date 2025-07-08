package handler

import (
	"auth/internal/model"
	"auth/internal/query"
	"auth/internal/setup"
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getJwtToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(setup.Env.SecretKey))
	return encodedToken, err
}

func Register(ctx context.Context, body *model.User) (*model.AuthResponse, *HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}
	if len(body.Password) < 8 {
		return nil, NewErr(http.StatusBadRequest, "Password is too short")
	}
	if query.IsLoginInUse(ctx, body.Login) {
		return nil, NewErr(http.StatusBadRequest, "Login is already in use")
	}

	passwordHash, err := hashPassword(body.Password)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	user := &model.User{Login: body.Login, Password: passwordHash}
	user, err = query.CreateUser(ctx, user)
	if err != nil {
		return nil, NewErr(http.StatusBadRequest, err.Error())
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{User: user, Token: encodedToken}, nil
}

func Login(ctx context.Context, body *model.User) (*model.AuthResponse, *HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}

	user, err := query.GetUserByLogin(ctx, body.Login)
	if err != nil {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	}
	if valid := checkPasswordHash(body.Password, user.Password); !valid {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{User: user, Token: encodedToken}, nil
}