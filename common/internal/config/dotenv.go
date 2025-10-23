package config

import (
	"common/pkg/helper"
	"common/pkg/sugar"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var Env envFields

func LoadEnv(workDir string) error {
	envPath := filepath.Join(workDir, "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("Error loading .env file:\n\t%w", err)
	}

	Env = envFields{
		VirtualRuntime: sugar.Or(helper.VirtualRuntime(os.Getenv("VIRTUAL_RUNTIME")), "localhost"),
		GrpcUrls: grpcClientUrls{
			Auth:  "localhost:" + os.Getenv("AUTH_PORT"),
			Order: "localhost:" + os.Getenv("ORDER_PORT"),
			Stock: "localhost:" + os.Getenv("STOCK_PORT"),
		},
		KafkaBrokerHosts:   []string{"localhost:" + os.Getenv("KAFKA_PORT")},
		KafkaSerialization: os.Getenv("KAFKA_SERIALIZATION"),
		AdminKey:           os.Getenv("AUTH_ADMIN_KEY"),
	}

	return nil
}

const AppName = "STOCK"

type envFields struct {
	VirtualRuntime     helper.VirtualRuntime
	GrpcUrls           grpcClientUrls
	KafkaBrokerHosts   []string
	KafkaSerialization string
	AdminKey           string
}

type grpcClientUrls struct {
	Auth  string
	Order string
	Stock string
}