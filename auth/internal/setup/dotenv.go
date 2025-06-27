package setup

import (
	"base/pkg/sugar"
	"fmt"
	"os"

	"strconv"

	"github.com/joho/godotenv"
)

var Env envConfig

func LoadEnvs(envPath string) error {
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("Error loading .env file:\n\t%w", err)
	}

	Env = envConfig{
		Port: sugar.Default(strconv.Atoi(getenv("PORT"))),
		Db: dbConfig{
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

type dbConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type envConfig struct {
	Port      int
	Db        dbConfig
	SecretKey string
	Test      bool
}

func getenv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", ServiceName, varName))
}