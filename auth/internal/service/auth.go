package service

import (
	"auth/internal/config"
	"auth/internal/model"
	"auth/internal/store"
	errpkg "common/pkg/errors"
	"common/pkg/grpcutil"
	commService "common/pkg/service"
	"common/pkg/sugar"
	"context"
	"errors"
	"net/http"
)

var NewErr = grpcutil.NewError
var NewInternalErr = grpcutil.NewInternalError

type userStore interface {
	Get(ctx context.Context, login string) (model.User, error)
	Save(ctx context.Context, user model.User) (model.User, error)
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
	body model.User,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}
	if len(body.Password) < 8 {
		return nil, NewErr(http.StatusBadRequest, "Password is too short")
	}

	if _, err := s.store.Get(ctx, body.Login); err == nil {
		return nil, NewErr(http.StatusBadRequest, "Login is already in use")
	} else if !errors.Is(err, errpkg.NotFound) {
		return nil, NewInternalErr(err)
	}

	user, err := s.store.Save(
		ctx,
		model.User{Login: body.Login, Password: body.Password, Role: model.UserRoles.Client},
	)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	token, err := createJwtToken(user)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	return &model.AuthResponse{Token: token}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	body model.User,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	if body.Login == "" || body.Password == "" {
		return nil, NewErr(http.StatusBadRequest, "Login and password must be provided")
	}

	user, err := s.store.Get(ctx, body.Login)
	if errors.Is(err, errpkg.NotFound) {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	} else if err != nil {
		return nil, NewInternalErr(err)
	}
	if valid := checkPasswordHash(body.Password, user.PasswordHash); !valid {
		return nil, NewErr(http.StatusUnauthorized, "Login or password is not valid")
	}

	token, err := createJwtToken(user)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	return &model.AuthResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(
	ctx context.Context,
	token string,
) *grpcutil.HandlerError {
	err := validateJwtToken(token)
	if err != nil {
		return NewErr(http.StatusUnauthorized, "Invalid Token")
	}
	return nil
}

func (s *AuthService) UpdateUser(
	ctx context.Context,
	token string,
	body model.User,
	roleKey string,
) (*model.AuthResponse, *grpcutil.HandlerError) {
	err := validateJwtToken(token)
	if err != nil {
		return nil, NewErr(http.StatusUnauthorized, "Invalid Token")
	}

	claims, err := commService.ParseJwtToken(token)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	oldUser, err := s.store.Get(ctx, claims.Login)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	updUser := model.User{
		Login:    sugar.If(body.Login != "", body.Login, oldUser.Login),
		Password: sugar.If(body.Password != "", body.Password, oldUser.Password),
		Role:     sugar.If(roleKey == config.Env.AdminKey, model.UserRoles.Admin, oldUser.Role),
	}

	updUser, err = s.store.Save(ctx, updUser)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	token, err = createJwtToken(updUser)
	if err != nil {
		return nil, NewInternalErr(err)
	}

	return &model.AuthResponse{Token: token}, nil
}