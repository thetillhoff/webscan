package webscan

import "github.com/thetillhoff/webscan/pkg/portScan"

func (engine Engine) ScanPortSimple() (Engine, error) {
	var (
		scanPorts = []uint16{
			80,  // HTTP
			443, // HTTPS
		}

		openPortsPerIp map[string][]uint16
	)

	openPortsPerIp = portScan.ScanPortRangeOfIps(append(engine.DnsScanEngine.ARecords, engine.dnsScanResult.AAAARecords...), scanPorts)

	engine.portScanOpenPorts, engine.portScanInconsistencies = portScan.CompareOpenPortsOfIps(openPortsPerIp)

	// Check if HTTP / HTTPS are available
	for _, openPort := range engine.portScanOpenPorts {
		if openPort == 80 {
			engine.isAvailableViaHttp = true
		} else if openPort == 443 {
			engine.isAvailableViaHttps = true
		}
	}

	return engine, nil
}
