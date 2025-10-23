package config

import (
	"common/pkg/helper"
	"common/pkg/sugar"
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
		Port:           sugar.Default(strconv.Atoi(getAppEnv("PORT"))),
		VirtualRuntime: sugar.Or(helper.VirtualRuntime(os.Getenv("VIRTUAL_RUNTIME")), "localhost"),
		TestPort:       sugar.Default(strconv.Atoi(getAppEnv("TEST_PORT"))),
		Test:           sugar.Default(strconv.ParseBool(getAppEnv("TEST"))),
		Db: dbConfig{
			User:        getAppEnv("DB_USER"),
			Password:    getAppEnv("DB_PASSWORD"),
			Port:        sugar.Default(strconv.Atoi(getAppEnv("DB_PORT"))),
			Name:        getAppEnv("DB_NAME"),
			SchemaReset: sugar.Default(strconv.ParseBool(getAppEnv("DB_SCHEMA_RESET"))),
		},
		GrpcUrls: grpcClientUrls{
			Auth:  "localhost:" + os.Getenv("AUTH_PORT"),
			Stock: "localhost:" + os.Getenv("STOCK_PORT"),
		},
		KafkaBrokerHosts:   []string{"localhost:" + os.Getenv("KAFKA_PORT")},
		KafkaSerialization: os.Getenv("KAFKA_SERIALIZATION"),
	}

	if Env.Test {
		Env.Db.Name += "_test"
	}

	return nil
}

const AppName = "ORDER"

type envFields struct {
	Port               int
	VirtualRuntime     helper.VirtualRuntime
	TestPort           int
	Test               bool
	Db                 dbConfig
	GrpcUrls           grpcClientUrls
	KafkaBrokerHosts   []string
	KafkaSerialization string
}

type dbConfig struct {
	User        string
	Password    string
	Port        int
	Name        string
	SchemaReset bool
}

type grpcClientUrls struct {
	Auth  string
	Stock string
}

func getAppEnv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}