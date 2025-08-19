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
		Port: sugar.Default(strconv.Atoi(getAppEnv("PORT"))),
		TestPort: sugar.Default(strconv.Atoi(getAppEnv("TEST_PORT"))),
		Test: sugar.Default(strconv.ParseBool(getAppEnv("TEST"))),
		Db: dbConfig{
			User:        getAppEnv("DB_USER"),
			Password:    getAppEnv("DB_PASSWORD"),
			Host:        getAppEnv("DB_HOST"),
			Port:        sugar.Default(strconv.Atoi(getAppEnv("DB_PORT"))),
			Name:        getAppEnv("DB_NAME"),
			SchemaReset: sugar.Default(strconv.ParseBool(getAppEnv("DB_SCHEMA_RESET"))),
		},
	}

	if Env.Test {
		Env.Db.Name += "_test"
	}

	return nil
}

const AppName = "STOCK"

type envFields struct {
	Port     int
	TestPort int
	Test     bool
	Db       dbConfig
}

type dbConfig struct {
	User        string
	Password    string
	Host        string
	Port        int
	Name        string
	SchemaReset bool
}

func getAppEnv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}