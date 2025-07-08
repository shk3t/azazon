package setup

import (
	"auth/internal/database"
	"base/pkg/log"
	baseSetup "base/pkg/setup"
	"os"
)

func initAll(envPath string, workDir string) error {
	if err := LoadEnv(envPath); err != nil {
		return err
	}
	if err := log.Init(workDir); err != nil {
		return err
	}
	if err := database.ConnectDatabase(Env.Db, Env.Test); err != nil {
		return err
	}

	if log.DLog != nil {
		log.DLog("Config inited successfully")
	}

	return nil
}

func deinitAll() {
	if log.DLog != nil {
		log.DLog("Config deinitialization...")
	}
	log.Deinit()
	database.ConnPool.Close()
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