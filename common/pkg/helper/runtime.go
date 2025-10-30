package helper

import (
	"fmt"
	"net"
)

type VirtualRuntime string

var VirtualRuntimes = struct {
	Localhost  VirtualRuntime
	Kubernetes VirtualRuntime
}{
	Localhost:  "localhost",
	Kubernetes: "kubernetes",
}

type OpMode string

var OpModes = struct {
	Read  OpMode
	Write OpMode
}{
	Read:  "read",
	Write: "write",
}

func (vr VirtualRuntime) GetDbHosts(appName string, mode OpMode) []string {
	switch vr {
	case VirtualRuntimes.Localhost:
		return []string{"localhost"}
	case VirtualRuntimes.Kubernetes:
		switch mode {
		case OpModes.Read:
			ips, err := net.LookupIP(fmt.Sprintf("%s-database-service", appName))
			if err != nil {
				panic(fmt.Errorf("DNS lookup failed: %v", err))
			}

			hosts := make([]string, len(ips))
			for i, ip := range ips {
				hosts[i] = ip.String()
			}
			return hosts
		case OpModes.Write:
			return []string{
				fmt.Sprintf(
					"%s-database-statefulset-0.%s-database-service.%s.svc.cluster.local",
					appName, appName, "azazon",
				),
			}
		default:
			panic("Unexpected operation mode")
		}
	default:
		panic("Unexpected virtual runtime")
	}
}