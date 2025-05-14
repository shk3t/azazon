package config

import (
	"auth/pkg/sugar"
	"fmt"
	"os"

	"strconv"

	"github.com/joho/godotenv"
)

var Env envConfig

func LoadEnvs() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file")
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
	}
}

const serviceName = "AUTH_"

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
}

func getenv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", serviceName, varName))
}