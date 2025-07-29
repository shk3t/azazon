package database

import (
	"auth/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var tableDefinitions = [...]string{
	`
    CREATE TABLE IF NOT EXISTS "user" (
        id SERIAL PRIMARY KEY,
        login VARCHAR(64) NOT NULL UNIQUE,
        password_hash VARCHAR(128) NOT NULL
    )`,
}

func InitDatabaseSchema(ctx context.Context) error {
	tx, _ := ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if config.Env.Test {
		tx.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		tx.Exec(ctx, "CREATE SCHEMA public")
	}

	for _, tableDef := range tableDefinitions {
		_, err := tx.Exec(ctx, tableDef)
		if err != nil {
			return fmt.Errorf("Schema initiation failed: %w", err)
		}
	}

	tx.Commit(ctx)
	return nil
}