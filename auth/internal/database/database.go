package database

import (
	"auth/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func Connect(ctx context.Context) {
	envDb := config.Env.Db

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		envDb.User, envDb.Password, envDb.Host, envDb.Port, envDb.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		panic(err)
	}

	InitSchema(ctx)
}