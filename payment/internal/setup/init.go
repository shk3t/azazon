package setup

import (
	"common/pkg/log"
	commSetup "common/pkg/setup"
	"payment/internal/config"
	"payment/internal/server"
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

var InitAll, DeinitAll = commSetup.CreateInitFuncs(initAll, deinitAll)