package portScan

import (
	"log/slog"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/status"
)

// maxConcurrentPortChecks caps in-flight TCP dials. Without it, a full
// 65535-port scan spawns one goroutine per port and exhausts file descriptors.
// ponytail: hardcoded cap, add a setting if a use case ever needs tuning it.
const maxConcurrentPortChecks = 500

func scanPortRangeOfIps(status *status.Status, ips []string, ports []uint16, timeout time.Duration) map[string][]uint16 {
	var (
		wg                 sync.WaitGroup
		openPortsPerIp     = map[string][]uint16{}
		ipPortTupleChannel = make(chan IpPortTuple, len(ips)*len(ports))
		sem                = make(chan struct{}, maxConcurrentPortChecks)
	)

	slog.Debug("portScan: Scanning port range of ips started", "len(ips)", len(ips), "len(ports)", len(ports))

	status.SpinningXOfInit(len(ips)*len(ports), "Scanning ports...")

	for _, ip := range ips { // For each ip
		for _, port := range ports { // For each port
			wg.Add(1)
			sem <- struct{}{} // Acquire a slot, blocking once maxConcurrentPortChecks are in flight
			go func(ip string, port uint16) {
				defer func() { <-sem }() // Release the slot
				isOpenTcpPort(
					&wg,
					status,
					IpPortTuple{
						Ip:   ip,
						Port: port,
					},
					ipPortTupleChannel,
					timeout,
				)
			}(ip, port)
		}
	}

	wg.Wait()                 // Wait until all goroutines are finished
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
