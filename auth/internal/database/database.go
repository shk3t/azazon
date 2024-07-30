package database

import (
	"auth/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func Connect(ctx context.Context) {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Env.Db.User,
		config.Env.Db.Password,
		config.Env.Db.Host,
		config.Env.Db.Port,
		config.Env.Db.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		panic(err)
	}

	InitSchema(ctx)
}