package setup

import (
	"common/pkg/log"
	setuppkg "common/pkg/setup"
	"stock/internal/config"
	"stock/internal/database"
	"stock/internal/server"
)

func initAll(workDir string) error {
	if err := config.LoadEnv(workDir); err != nil {
		return err
	}
	if err := log.Init(workDir); err != nil {
		return err
	}
	if err := database.ConnectDatabase(workDir); err != nil {
		return err
	}

	return nil
}

func deinitAll() {
	server.Deinit()
	database.ConnPool.Close()
	log.Deinit()
}

var InitAll, DeinitAll = setuppkg.CreateInitFuncs(initAll, deinitAll)