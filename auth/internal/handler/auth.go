package handler

import (
	"auth/internal/config"
	m "auth/internal/model"
	q "auth/internal/query"
	"context"
	"encoding/json"
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

func getJwtToken(user *m.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(config.Env.SecretKey))
	return encodedToken, err
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := m.User{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Login == "" || body.Password == "" {
		http.Error(w, "Login and password must be provided", http.StatusBadRequest)
		return
	}
	if len(body.Password) < 8 {
		http.Error(w, "Password is too short", http.StatusBadRequest)
		return
	}
	if q.IsLoginInUse(ctx, body.Login) {
		http.Error(w, "Login is already in use", http.StatusBadRequest)
		return
	}

	passwordHash, err := hashPassword(body.Password)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user := &m.User{Login: body.Login, Password: passwordHash}
	user, err = q.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m.AuthResponse{User: user, Token: encodedToken})
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body := m.User{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Login == "" || body.Password == "" {
		http.Error(w, "Login and password must be provided", http.StatusBadRequest)
		return
	}

	user, err := q.GetUserByLogin(ctx, body.Login)
	if err != nil {
		http.Error(w, "Login or password is not valid", http.StatusUnauthorized)
		return
	}
	if valid := checkPasswordHash(body.Password, user.Password); !valid {
		http.Error(w, "Login or password is not valid", http.StatusUnauthorized)
		return
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m.AuthResponse{User: user, Token: encodedToken})
}