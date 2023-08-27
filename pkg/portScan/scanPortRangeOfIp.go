package portScan

import (
	"fmt"
	"net"
	"sync"
)

var wgPortScan sync.WaitGroup

func scanPortRangeOfIp(ip string, ports []uint16, verbose bool, portChannel chan<- IpPortsTuple) {
	var (
		openPorts       = []uint16{}
		openPortChannel = make(chan uint16, len(ports))
	)
	defer wgIpScan.Done()

	for _, port := range ports { // For each important port
		wgPortScan.Add(1)                                        // Wait for one more goroutine to finish
		go isOpenTcpPort(net.ParseIP(ip), port, openPortChannel) // Start goroutine that checks if port is open
		if verbose {
			fmt.Println("Started scanning port", port, "of ip", ip)
		}
	}

	wgPortScan.Wait() // Wait until all goroutines are finished

	close(openPortChannel) // Make sure channel is closed when goroutines are finished

	for openPort := range openPortChannel { // Convert from channel to slice
		openPorts = append(openPorts, openPort)
	}

	portChannel <- IpPortsTuple{ // Return equivalent
		Ip:    ip,
		Ports: openPorts,
	}
}
