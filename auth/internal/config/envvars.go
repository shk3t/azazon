package config

import (
	"auth/pkg/sugar"
	"os"

	"strconv"

	"github.com/joho/godotenv"
)

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

var Env envConfig

func LoadEnvs() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}

	Env = envConfig{
		Port: sugar.Default(strconv.Atoi(os.Getenv("PORT"))),
		Db: dbConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     sugar.Default(strconv.Atoi(os.Getenv("DB_PORT"))),
			Name:     os.Getenv("DB_NAME"),
		},
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}