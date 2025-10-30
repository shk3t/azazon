package database

import (
	"context"
)

var Pooler *dbPooler

func ConnectDatabase(workDir string) error {
	var err error
	ctx := context.Background()
	Pooler, err = NewDbPooler(ctx)
	if err != nil {
		return err
	}
	return InitDatabaseSchema(ctx, workDir)
}