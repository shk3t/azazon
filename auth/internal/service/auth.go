package service

import (
	"auth/internal/setup"
	"auth/internal/store"
	errorpkg "base/pkg/error"
	"base/pkg/grpcutil"
	"base/pkg/model"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var NewErr = grpcutil.NewError

type userStore interface {
	Get(ctx context.Context, login string) (*model.User, error)
	Save(ctx context.Context, user model.User) (*model.User, error)
}

type AuthService struct {
	store userStore
}

func NewAuthService() *AuthService {
	return &AuthService{
		store: &store.PostgreUserStore{
			HashPassword: hashPassword,
		},
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	body *model.User,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}
	if len(body.Password) < 8 {
		return nil, NewErr(http.StatusBadRequest, "Password is too short")
	}

	if user, err := s.store.Get(ctx, body.Login); user != nil {
		return nil, NewErr(http.StatusBadRequest, "Login is already in use")
	} else if !errors.Is(err, errorpkg.NotFound) {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	user, err := s.store.Save(ctx, *body)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{User: user, Token: encodedToken}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	body *model.User,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}

	user, err := s.store.Get(ctx, body.Login)
	if errors.Is(err, errorpkg.NotFound) {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	} else if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}
	if valid := checkPasswordHash(body.Login, user.PasswordHash); !valid {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{User: user, Token: encodedToken}, nil
}

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
		"exp":   time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(setup.Env.SecretKey))
	return encodedToken, err
}