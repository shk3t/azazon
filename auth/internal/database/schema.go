package database

import (
	"auth/internal/config"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var tableDefinitions = [...]string{
	`
    CREATE TABLE IF NOT EXISTS "user" (
        id SERIAL PRIMARY KEY,
        login VARCHAR(64) NOT NULL UNIQUE,
        password VARCHAR(128) NOT NULL
    )`,
}

func InitSchema(ctx context.Context) {
	tx, _ := ConnPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	if config.Env.Test {
		tx.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		tx.Exec(ctx, "CREATE SCHEMA public")
	}

	for _, tableDef := range tableDefinitions {
		_, err := tx.Exec(ctx, tableDef)
		if err != nil {
			panic("Schema initiation failed: " + err.Error())
		}
	}

	tx.Commit(ctx)
	log.Println("Schema inited successfully!")
}