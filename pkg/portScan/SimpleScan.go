package portScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

func SimpleScan(status *status.Status, aRecords []string, aaaaRecords []string) (Result, error) {
	var (
		scanPorts = []uint16{
			80,  // HTTP
			443, // HTTPS
		}

		result = Result{
			openPorts:               []uint16{},
			openPortInconsistencies: []string{},
			isAvailableViaHttp:      false,
			isAvailableViaHttps:     false,
		}

		openPortsPerIp map[string][]uint16
	)

	slog.Debug("portScan: Simple scan started")

	openPortsPerIp = scanPortRangeOfIps(status, append(aRecords, aaaaRecords...), scanPorts)

	result.openPorts, result.openPortInconsistencies = CompareOpenPortsOfIps(openPortsPerIp)

	// Check if HTTP / HTTPS are available
	for _, openPort := range result.openPorts {
		if openPort == 80 {
			result.isAvailableViaHttp = true
		} else if openPort == 443 {
			result.isAvailableViaHttps = true
		}
	}

	slog.Debug("portScan: Advanced scan completed")

	return result, nil
}
