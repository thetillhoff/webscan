package portScan

import (
	"log/slog"
	"slices"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func SimpleScan(target types.Target, status *status.Status, aRecords []string, aaaaRecords []string) (Result, error) {
	var (
		result = Result{
			OpenPortsPerIp:          map[string][]uint16{},
			openPorts:               []uint16{},
			openPortInconsistencies: []string{},
			isPort80Open:            false,
			isPort443Open:           false,
		}

		scanPorts []uint16
	)

	slog.Debug("portScan: Simple scan started", "len(aRecords)", len(aRecords), "len(aaaaRecords)", len(aaaaRecords))

	switch { // Selecting which ports to scan
	case target.Port() != "":
		scanPorts = []uint16{target.PortAsUint16()}
	case target.Schema() == types.HTTPS:
		scanPorts = []uint16{443}
	case target.Schema() == types.HTTP:
		scanPorts = []uint16{80}
	default:
		scanPorts = []uint16{80, 443}
	}

	result.OpenPortsPerIp = scanPortRangeOfIps(status, append(aRecords, aaaaRecords...), scanPorts)

	result.openPorts, result.openPortInconsistencies = CompareOpenPortsOfIps(result.OpenPortsPerIp)

	// Check if port 80 / 443 are open
	result.isPort80Open = slices.Contains(result.openPorts, 80)
	result.isPort443Open = slices.Contains(result.openPorts, 443)

	slog.Debug("portScan: Simple scan completed", "len(result.openPorts)", len(result.openPorts))

	return result, nil
}
