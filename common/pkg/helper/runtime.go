package helper

import "fmt"

type VirtualRuntime string

var VirtualRuntimes = struct {
	Localhost  VirtualRuntime
	Kubernetes VirtualRuntime
}{
	Localhost:  "localhost",
	Kubernetes: "kubernetes",
}

func (vr VirtualRuntime) GetDbHost(appName string) string {
	switch vr {
	case VirtualRuntimes.Localhost:
		return "localhost"
	case VirtualRuntimes.Kubernetes:
		return fmt.Sprintf("%s-database-service", appName)
	default:
		panic("Unexpected runtimeHost")
	}
}