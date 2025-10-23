package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type abstractConnection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(
		ctx context.Context,
		tableName pgx.Identifier,
		columnNames []string,
		rowSrc pgx.CopyFromSource,
	) (int64, error)
}

func TxOrPool(tx pgx.Tx, pool *pgxpool.Pool) abstractConnection {
	if tx != nil {
		return tx
	}
	return pool
}

func RunPgxSqlScript(conn *pgxpool.Pool, scriptPath string) error {
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("Failed to read script file: %w", err)
	}

	statements := strings.Split(string(scriptContent), ";")

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		stmt += ";"

		_, err := conn.Exec(context.Background(), stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement %d: %w\nStatement: %s", i+1, err, stmt)
		}
	}

	return nil
}