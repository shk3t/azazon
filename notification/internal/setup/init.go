package setup

import (
	"common/pkg/log"
	commSetup "common/pkg/setup"
	"notification/internal/config"
	"notification/internal/server"
	"notification/internal/service"
)

func initAll(workDir string) error {
	if err := config.LoadEnv(workDir); err != nil {
		return err
	}
	if err := log.Init(workDir); err != nil {
		return err
	}
	if err := service.InitMailer(workDir); err != nil {
		return err
	}

	return nil
}

func deinitAll() {
	server.Deinit()
	log.Deinit()
	service.DeinitMailer()
}

var InitAll, DeinitAll = commSetup.CreateInitFuncs(initAll, deinitAll)