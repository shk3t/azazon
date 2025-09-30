package config

import (
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

	Env = envFields{
		Port:     sugar.Default(strconv.Atoi(getAppEnv("PORT"))),
		TestPort: sugar.Default(strconv.Atoi(getAppEnv("TEST_PORT"))),
		Test:     sugar.Default(strconv.ParseBool(getAppEnv("TEST"))),
		PayTimeout: time.Second * time.Duration(
			sugar.Default(strconv.Atoi(getAppEnv("PAY_TIMEOUT"))),
		),
		KafkaBrokerHosts:   []string{"localhost:" + os.Getenv("KAFKA_PORT")},
		KafkaSerialization: os.Getenv("KAFKA_SERIALIZATION"),
	}

	return nil
}

const AppName = "PAYMENT"

type envFields struct {
	Port               int
	TestPort           int
	Test               bool
	PayTimeout         time.Duration
	KafkaBrokerHosts   []string
	KafkaSerialization string
}

func getAppEnv(varName string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", AppName, varName))
}