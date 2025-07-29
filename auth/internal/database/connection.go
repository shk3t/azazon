package database

import (
	"auth/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func ConnectDatabase() error {
	ctx := context.Background()
	db := config.Env.Db

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, db.Host, db.Port, db.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return err
	}

	return InitDatabaseSchema(ctx)
}