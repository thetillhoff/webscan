package portScan

import (
	"net"
	"sync"
)

var wg sync.WaitGroup

func ScanPortRangeOfIp(ip string, ports []uint16) []uint16 {
	var (
		openPorts       = []uint16{}
		openPortChannel = make(chan uint16, len(ports))
	)

	for _, port := range ports { // For each important port
		wg.Add(1)                                                // Wait for one more goroutine to finish
		go isOpenTcpPort(net.ParseIP(ip), port, openPortChannel) // Start goroutine that checks if port is open
	}

	wg.Wait() // Wait until all goroutines are finished

	close(openPortChannel) // Make sure channel is closed when goroutines are finished

	for openPort := range openPortChannel { // Convert from channel to slice
		openPorts = append(openPorts, openPort)
	}

	return openPorts
}
