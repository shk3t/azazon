package database

import (
	"auth/internal/config"
	"common/pkg/helper"
	"context"
	"fmt"
	"path/filepath"
)

func InitDatabaseSchema(ctx context.Context, workDir string) error {
	if !(config.Env.Db.SchemaReset || config.Env.Test) {
		return nil
	}

	Pooler.Writer().Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
	Pooler.Writer().Exec(ctx, "CREATE SCHEMA public")

	schemaFile := filepath.Join(workDir, "migrations", "init.sql")
	err := helper.RunPgxSqlScript(Pooler.Writer(), schemaFile)
	if err != nil {
		return fmt.Errorf("Schema initialization failed: %w", err)
	}

	return nil
}