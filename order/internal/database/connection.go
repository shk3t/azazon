package database

import (
	"context"
	"fmt"
	"order/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ConnPool *pgxpool.Pool

func ConnectDatabase(workDir string) error {
	ctx := context.Background()
	db := config.Env.Db

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User, db.Password,
		config.Env.VirtualRuntime.GetDbHost(config.AppName),
		db.Port, db.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, dbUrl)
	if err != nil {
		return err
	}

	return InitDatabaseSchema(ctx, workDir)
}