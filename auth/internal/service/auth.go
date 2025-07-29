package service

import (
	"auth/internal/store"
	errpkg "base/pkg/error"
	"base/pkg/grpcutil"
	"base/pkg/model"
	"context"
	"errors"
	"net/http"
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
	} else if !errors.Is(err, errpkg.NotFound) {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	user, err := s.store.Save(ctx, *body)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	token, err := encodeJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{Token: token}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	body *model.User,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}

	user, err := s.store.Get(ctx, body.Login)
	if errors.Is(err, errpkg.NotFound) {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	} else if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}
	if valid := checkPasswordHash(body.Login, user.PasswordHash); !valid {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	}

	token, err := encodeJwtToken(user)
	if err != nil {
		return nil, NewErr(http.StatusInternalServerError, "")
	}

	return &model.AuthResponse{Token: token}, nil
}