package database

import (
	"order/internal/config"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

func InitDatabaseSchema(ctx context.Context, workDir string, dbUrl string) error {
	if !(config.Env.Db.SchemaReset || config.Env.Test) {
		return nil
	}

	ConnPool.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
	ConnPool.Exec(ctx, "CREATE SCHEMA public")

	schemaFile := filepath.Join(workDir, "migrations", "init.sql")
	err := runPsqlScript(dbUrl, schemaFile)
	if err != nil {
		return fmt.Errorf("Schema initialization failed: %w", err)
	}

	return nil
}

func runPsqlScript(connString, scriptPath string) error {
	cmd := exec.Command("psql", connString, "-f", scriptPath)

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stderrStr := stderrBuf.String()
	if stderrStr != "" {
		return errors.New(stderrStr)
	}

	return err
}