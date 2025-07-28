package store

import (
	"auth/internal/query"
	errorpkg "base/pkg/error"
	"base/pkg/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type PostgreUserStore struct {
	HashPassword func(password string) (string, error)
}

func (s *PostgreUserStore) Get(ctx context.Context, login string) (*model.User, error) {
	user, err := query.GetUserByLogin(ctx, login)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorpkg.NotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgreUserStore) Save(ctx context.Context, user model.User) (*model.User, error) {
	var err error

	user.PasswordHash, err = s.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	if user.Id == 0 {
		user.Id, err = query.CreateUser(ctx, user)
	} else {
		err = query.UpdateUser(ctx, user.Id, user)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}