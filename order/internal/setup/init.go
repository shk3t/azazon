package setup

import (
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"order/internal/config"
	"order/internal/database"
	"order/internal/server"
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

var InitAll, DeinitAll = baseSetup.CreateInitFuncs(initAll, deinitAll)