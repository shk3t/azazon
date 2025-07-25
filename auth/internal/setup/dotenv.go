package setup

import (
	"auth/internal/database"
	"base/pkg/sugar"
	"fmt"
	"os"

	"strconv"

	"github.com/joho/godotenv"
)

var Env envFields

func LoadEnv(envPath string) error {
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("Error loading .env file:\n\t%w", err)
	}

	Env = envFields{
		Port: sugar.Default(strconv.Atoi(getenv("PORT"))),
		Db: database.DbConfig{
			User:     getenv("DB_USER"),
			Password: getenv("DB_PASSWORD"),
			Host:     getenv("DB_HOST"),
			Port:     sugar.Default(strconv.Atoi(getenv("DB_PORT"))),
			Name:     getenv("DB_NAME"),
		},
		SecretKey: getenv("SECRET_KEY"),
		Test:      sugar.Default(strconv.ParseBool(getenv("TEST"))),
	}

	if Env.Test {
		Env.Db.Name += "_test"
	}

	return nil
}

const ServiceName = "AUTH"

type envFields struct {
	Port      int
	Db        database.DbConfig
	SecretKey string
	Test      bool
}

func getenv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", ServiceName, varName))
}