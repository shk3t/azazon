package setup

import (
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"notification/internal/config"
	"notification/internal/server"
)

func initAll(workDir string) error {
	if err := config.LoadEnv(workDir); err != nil {
		return err
	}
	if err := log.Init(workDir); err != nil {
		return err
	}

	return nil
}

func deinitAll() {
	server.Deinit()
	log.Deinit()
}

var InitAll, DeinitAll = baseSetup.CreateInitFuncs(initAll, deinitAll)