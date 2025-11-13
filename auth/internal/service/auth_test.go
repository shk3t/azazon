package service

import (
	"auth/internal/model"
	errpkg "common/pkg/errors"
	"common/pkg/sugar"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		ctx  context.Context
		body model.User
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				body: model.User{
					Login:    "man",
					Password: "somepassword",
				},
			},
			err: nil,
		},
		{
			name: "Short password",
			args: args{
				ctx: context.Background(),
				body: model.User{
					Login:    "shortman",
					Password: "short",
				},
			},
			err: NewErr(http.StatusBadRequest, "Password is too short"),
		},
	}

	mockStore := newMockuserStore(t)
	mockStore.EXPECT().Get(
		mock.Anything,
		mock.AnythingOfType("string"),
	).Return(model.User{}, errpkg.NotFound)
	mockStore.EXPECT().Save(
		mock.Anything,
		mock.AnythingOfType("model.User"),
	).RunAndReturn(
		func(ctx context.Context, user model.User) (model.User, error) {
			return user, nil
		},
	)
	s := &AuthService{store: mockStore}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Register(tt.args.ctx, tt.args.body)
			if tt.err == nil {
				assert.NotNil(resp)
			}
			assert.Equal(tt.err, err.Interface())
		})
	}
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		ctx  context.Context
		body model.User
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Ok",
			args: args{
				ctx: context.Background(),
				body: model.User{
					Login:    "man",
					Password: "somepassword",
				},
			},
			err: nil,
		},
		{
			name: "Login is not provided",
			args: args{
				ctx: context.Background(),
				body: model.User{
					Login:    "",
					Password: "somepassword",
				},
			},
			err: NewErr(http.StatusBadRequest, "Login and password must be provided"),
		},
	}

	mockStore := newMockuserStore(t)
	mockStore.EXPECT().Get(
		mock.Anything, "man",
	).Return(
		model.User{
			Login: "man",
			PasswordHash: sugar.Default(hashPassword("somepassword")),
		},
		nil,
	)
	s := &AuthService{store: mockStore}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := s.Login(tt.args.ctx, tt.args.body)
			if tt.err == nil {
				assert.NotNil(resp)
			}
			assert.Equal(tt.err, err.Interface())
		})
	}
}