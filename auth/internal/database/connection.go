package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

var ConnPool *pgxpool.Pool

func ConnectDatabase(cfg DbConfig, isTestEnv bool) error {
	ctx := context.Background()

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	var err error
	ConnPool, err = pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return err
	}

	return InitDatabaseSchema(ctx, isTestEnv)
}