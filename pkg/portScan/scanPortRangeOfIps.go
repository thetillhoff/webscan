package portScan

import (
	"log/slog"
	"sync"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

// var wgIpScan sync.WaitGroup
var wgPortScan sync.WaitGroup

func scanPortRangeOfIps(status *status.Status, ips []string, ports []uint16) map[string][]uint16 {
	var (
		openPortsPerIp     = map[string][]uint16{}
		ipPortTupleChannel = make(chan IpPortTuple, len(ips)*len(ports))
	)

	slog.Debug("portScan: Scanning port range of ips started", "len(ips)", len(ips), "len(ports)", len(ports))

	status.SpinningXOfInit(len(ips)*len(ports), "Scanning ports...")

	for _, ip := range ips { // For each ip
		for _, port := range ports { // For each port
			wgPortScan.Add(1)
			go isOpenTcpPort(
				status,
				IpPortTuple{
					Ip:   ip,
					Port: port,
				},
				ipPortTupleChannel,
			) // Start goroutine that checks if port is open
		}
	}

	wgPortScan.Wait()         // Wait until all goroutines are finished
	close(ipPortTupleChannel) // Make sure channel is closed when goroutines are finished
	status.SpinningXOfComplete("Scan of open ports completed.")

	for IpPortsTuple := range ipPortTupleChannel { // Convert from channel to slice
		if _, ok := openPortsPerIp[IpPortsTuple.Ip]; !ok { // Create list for ip if not exists
			openPortsPerIp[IpPortsTuple.Ip] = []uint16{} // Init list
		}

		openPortsPerIp[IpPortsTuple.Ip] = append(openPortsPerIp[IpPortsTuple.Ip], IpPortsTuple.Port)
	}

	slog.Debug("portScan: Scanning port range of ips completed")

	return openPortsPerIp
}
