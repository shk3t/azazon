package setup

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/server"
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"os"
)

func initAll(envPath string, workDir string) error {
	if err := config.LoadEnv(envPath); err != nil {
		return err
	}
	if err := log.Init(workDir); err != nil {
		return err
	}
	if err := database.ConnectDatabase(); err != nil {
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
		return initAll(args[0].(string), args[1].(string))
	},
	deinitAll,
)
var InitAll = func(envPath string, workDir string) error {
	return initializer.Init(envPath, workDir)
}
var DeinitAll = initializer.Deinit

func GracefullExit(code int) {
	DeinitAll()
	os.Exit(code)
}