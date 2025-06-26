package setup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func ConnectDatabase() error {
	envDb := Env.Db
	ctx := context.Background()

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		envDb.User, envDb.Password, envDb.Host, envDb.Port, envDb.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return err
	}

	return InitDatabaseSchema(ctx)
}