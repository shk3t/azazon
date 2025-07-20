package store

import (
	"context"
	"auth/internal/model"
)

type UserStore interface {
	Get(ctx context.Context, login string) (model.User, error)
	Save(ctx context.Context, login string, password string) (int, error)
}