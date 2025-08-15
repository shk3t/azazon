package config

import (
	"base/pkg/sugar"
	"fmt"
	"os"
	"path/filepath"

	"strconv"

	"github.com/joho/godotenv"
)

var Env envFields

func LoadEnv(workDir string) error {
	envPath := filepath.Join(workDir, "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("Error loading .env file:\n\t%w", err)
	}

	Env = envFields{
		Port: sugar.Default(strconv.Atoi(getenv("PORT"))),
		Test:      sugar.Default(strconv.ParseBool(getenv("TEST"))),
	}

	return nil
}

const AppName = "ORDER"

type envFields struct {
	Port      int
	Test      bool
}

func getenv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}