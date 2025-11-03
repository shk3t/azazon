package config

import (
	"common/pkg/helper"
	"common/pkg/sugar"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"strconv"

	"github.com/joho/godotenv"
)

var Env envFields

func LoadEnv(workDir string) error {
	envPath := filepath.Join(workDir, "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("Error loading .env file:\n\t%w", err)
	}

	virtualRuntime := sugar.If(
		os.Getenv("VIRTUAL_RUNTIME") == string(helper.VirtualRuntimes.Kubernetes),
		helper.VirtualRuntimes.Kubernetes,
		helper.VirtualRuntimes.Localhost,
	)

	Env = envFields{
		Port: sugar.Default(strconv.Atoi(getAppEnv("PORT"))),
		VirtualRuntime: virtualRuntime,
		TestPort: sugar.Default(strconv.Atoi(getAppEnv("TEST_PORT"))),
		Test:     sugar.Default(strconv.ParseBool(getAppEnv("TEST"))),
		Db: dbConfig{
			User:        getAppEnv("DB_USER"),
			Password:    getAppEnv("DB_PASSWORD"),
			Port:        sugar.Default(strconv.Atoi(getAppEnv("DB_PORT"))),
			Name:        getAppEnv("DB_NAME"),
			SchemaReset: sugar.Default(strconv.ParseBool(getAppEnv("DB_SCHEMA_RESET"))),
		},
		PayTimeout: time.Second * time.Duration(
			sugar.Default(strconv.Atoi(getAppEnv("PAY_TIMEOUT"))),
		),
		KafkaBrokerHosts:   virtualRuntime.GetKafkaHosts(),
		KafkaSerialization: os.Getenv("KAFKA_SERIALIZATION"),
	}

	return nil
}

const AppName = "PAYMENT"

type envFields struct {
	Port               int
	VirtualRuntime     helper.VirtualRuntime
	TestPort           int
	Test               bool
	Db                 dbConfig
	PayTimeout         time.Duration
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

func getAppEnv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}