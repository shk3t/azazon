package store

import (
	"auth/internal/conversion"
	"auth/internal/database"
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
	user, err := query.New(database.Pooler.Reader()).GetUserByLogin(ctx, login)

	if errors.Is(err, pgx.ErrNoRows) {
		return *conversion.UserModel(&user), errpkg.NotFound
	}

	return *conversion.UserModel(&user), err
}

func (s *PostgreUserStore) Save(ctx context.Context, user model.User) (model.User, error) {
	var err error
	q := query.New(database.Pooler.Writer())

	if user.PasswordHash == "" {
		user.PasswordHash, err = s.HashPassword(user.Password)
	}
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		var id int32
		id, err = q.CreateUser(ctx, query.CreateUserParams{
			Login:        user.Login,
			PasswordHash: user.PasswordHash,
			Role:         string(user.Role),
		})
		user.Id = int(id)
	} else {
		err = q.UpdateUser(ctx, query.UpdateUserParams{
			ID:           int32(user.Id),
			Login:        user.Login,
			PasswordHash: user.PasswordHash,
			Role:         string(user.Role),
		})
	}

	return user, err
}