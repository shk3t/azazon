package database

import (
	"auth/internal/config"
	"common/pkg/helper"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbPooler struct {
	readConnPools []*pgxpool.Pool
	writeConnPool *pgxpool.Pool
	readCalls     int
}

func NewDbPooler(ctx context.Context) (*dbPooler, error) {
	db := config.Env.Db

	readHosts := config.Env.VirtualRuntime.GetDbHosts(config.AppName, helper.OpModes.Read)
	rPools := make([]*pgxpool.Pool, len(readHosts))
	for i, host := range readHosts {
		rDbUrl := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			db.User, db.Password, host, db.Port, db.Name,
		)
		rCfg, err := pgxpool.ParseConfig(rDbUrl)
		if err != nil {
			return nil, err
		}
		rCfg.MinConns = 1
		rCfg.MaxConns = 3

		pool, err := pgxpool.NewWithConfig(ctx, rCfg)
		if err != nil {
			return nil, err
		}
		rPools[i] = pool
	}

	writeHost := config.Env.VirtualRuntime.GetDbHosts(config.AppName, helper.OpModes.Write)[0]
	wDbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, writeHost, db.Port, db.Name,
	)
	wCfg, err := pgxpool.ParseConfig(wDbUrl)
	if err != nil {
		return nil, err
	}
	wCfg.MinConns = 1
	wCfg.MaxConns = 5

	wPool, err := pgxpool.NewWithConfig(ctx, wCfg)
	if err != nil {
		return nil, err
	}

	return &dbPooler{
		readConnPools: rPools,
		writeConnPool: wPool,
		readCalls:     0,
	}, nil
}

func (p *dbPooler) Reader() *pgxpool.Pool {
	p.readCalls++
	return p.readConnPools[p.readCalls%len(p.readConnPools)]
}

func (p *dbPooler) Writer() *pgxpool.Pool {
	return p.writeConnPool
}

func (p *dbPooler) Close() {
	for _, rPool := range p.readConnPools {
		rPool.Close()
	}
	p.writeConnPool.Close()
}