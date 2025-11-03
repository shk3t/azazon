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

	virtualRuntime := sugar.If(
		os.Getenv("VIRTUAL_RUNTIME") == string(helper.VirtualRuntimes.Kubernetes),
		helper.VirtualRuntimes.Kubernetes,
		helper.VirtualRuntimes.Localhost,
	)
	externalClusterIp := os.Getenv("EXTERNAL_CLUSTER_IP")
	externalClusterPort := os.Getenv("EXTERNAL_CLUSTER_PORT")

	Env = envFields{
		Domain:             os.Getenv("DOMAIN"),
		VirtualRuntime:     virtualRuntime,
		KafkaBrokerHosts:   virtualRuntime.GetKafkaHosts(),
		KafkaSerialization: os.Getenv("KAFKA_SERIALIZATION"),
		AdminKey:           os.Getenv("AUTH_ADMIN_KEY"),
	}

	switch virtualRuntime {
	case helper.VirtualRuntimes.Localhost:
		Env.GrpcUrls = grpcClientUrls{
			Auth:  "localhost:" + os.Getenv("AUTH_PORT"),
			Order: "localhost:" + os.Getenv("ORDER_PORT"),
			Stock: "localhost:" + os.Getenv("STOCK_PORT"),
		}
	case helper.VirtualRuntimes.Kubernetes:
		url := externalClusterIp + sugar.If(externalClusterPort != "", ":"+externalClusterPort, "")
		Env.GrpcUrls = grpcClientUrls{Auth: url, Order: url, Stock: url}
	}

	return nil
}

const AppName = "STOCK"

type envFields struct {
	Domain             string
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