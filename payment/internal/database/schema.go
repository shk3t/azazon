package database

import (
	"common/pkg/helper"
	"context"
	"fmt"
	"path/filepath"
	"payment/internal/config"
)

func InitDatabaseSchema(ctx context.Context, workDir string) error {
	if !(config.Env.Db.SchemaReset || config.Env.Test) {
		return nil
	}

	ConnPool.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
	ConnPool.Exec(ctx, "CREATE SCHEMA public")

	schemaFile := filepath.Join(workDir, "migrations", "init.sql")
	err := helper.RunPgxSqlScript(ConnPool, schemaFile)
	if err != nil {
		return fmt.Errorf("Schema initialization failed: %w", err)
	}

	return nil
}