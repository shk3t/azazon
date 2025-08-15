package setup

import (
	"stock/internal/config"
	"stock/internal/database"
	"stock/internal/server"
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"os"
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

var initializer = baseSetup.NewInitializer(
	func(args ...any) error {
		return initAll(args[0].(string))
	},
	deinitAll,
)
var InitAll = func(workDir string) error {
	return initializer.Init(workDir)
}
var DeinitAll = initializer.Deinit

func GracefullExit(code int) {
	DeinitAll()
	os.Exit(code)
}