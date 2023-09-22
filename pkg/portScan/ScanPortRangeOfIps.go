package portScan

import (
	"fmt"
	"sync"
)

var wgIpScan sync.WaitGroup

type IpPortsTuple struct {
	Ip    string
	Ports []uint16
}

func ScanPortRangeOfIps(ips []string, ports []uint16, verbose bool) map[string][]uint16 {
	var (
		openPortsPerIp      = map[string][]uint16{}
		IpPortsTupleChannel = make(chan IpPortsTuple, len(ips))
	)

	for _, ip := range ips { // For each ip
		wgIpScan.Add(1)                                               // Wait for one more goroutine to finish
		go scanPortRangeOfIp(ip, ports, verbose, IpPortsTupleChannel) // Start goroutine that checks port range for ip
		if verbose {
			fmt.Println("Started scanning ports of ip", ip)
		}
	}

	wgIpScan.Wait() // Wait until all goroutines are finished

	close(IpPortsTupleChannel) // Make sure channel is closed when goroutines are finished

	for IpPortsTuple := range IpPortsTupleChannel { // Convert from channel to slice
		openPortsPerIp[IpPortsTuple.Ip] = IpPortsTuple.Ports
	}

	return openPortsPerIp
}
