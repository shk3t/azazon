package store

import (
	"auth/internal/model"
	"auth/internal/query"
	errpkg "common/pkg/errors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type PostgreUserStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreUserStore) Get(ctx context.Context, login string) (model.User, error) {
	user, err := query.GetUserByLogin(ctx, login)

	if errors.Is(err, pgx.ErrNoRows) {
		return user, errpkg.NotFound
	}

	return user, err
}

func (s *PostgreUserStore) Save(ctx context.Context, user model.User) (model.User, error) {
	var err error

	user.PasswordHash, err = s.HashPassword(user.Password)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		user.Id, err = query.CreateUser(ctx, user)
	} else {
		err = query.UpdateUser(ctx, user.Id, user)
	}

	return user, err
}