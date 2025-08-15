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
		Port:     sugar.Default(strconv.Atoi(getenv("PORT"))),
		TestPort: sugar.Default(strconv.Atoi(getenv("TEST_PORT"))),
		Test:     sugar.Default(strconv.ParseBool(getenv("TEST"))),
		Db: dbConfig{
			User:        getenv("DB_USER"),
			Password:    getenv("DB_PASSWORD"),
			Host:        getenv("DB_HOST"),
			Port:        sugar.Default(strconv.Atoi(getenv("DB_PORT"))),
			Name:        getenv("DB_NAME"),
			SchemaReset: sugar.Default(strconv.ParseBool(getenv("DB_SCHEMA_RESET"))),
		},
		SecretKey: getenv("SECRET_KEY"),
		AdminKey:  getenv("ADMIN_KEY"),
	}

	if Env.Test {
		Env.Db.Name += "_test"
	}

	return nil
}

const AppName = "AUTH"

type envFields struct {
	Port      int
	TestPort  int
	Test      bool
	Db        dbConfig
	SecretKey string
	AdminKey  string
}

type dbConfig struct {
	User        string
	Password    string
	Host        string
	Port        int
	Name        string
	SchemaReset bool
}

func getenv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}