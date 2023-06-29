package portScan

import (
	"strconv"
	"strings"
)

// Checks whether the open ports are the same of all ip addresses
// Returns list of ports that where open on all ip addresses,
// a list of inconsistencies
func CompareOpenPortsOfIps(openPortsPerIp map[string][]uint16) ([]uint16, []string) {
	var (
		uniqueRelevantPorts  = map[uint16]struct{}{}
		ipsWithPortOpened    []string
		ipsWithoutPortOpened []string

		openPorts       = []uint16{}
		inconsistencies = []string{}
	)

	// Create map where each uint16 value of all ips is an empty struct -> result is a list of unique relevant ports
	for _, openPorts := range openPortsPerIp { // Iterate over ips
		for _, openPort := range openPorts { // Iterate over openPorts of ips
			uniqueRelevantPorts[openPort] = struct{}{} // Add openPort to list of unique relevant ports
		}
	}

	for uniqueRelevantPort := range uniqueRelevantPorts { // Iterate over unique relevant ports and compare to actual open ports of each ip

		ipsWithPortOpened = []string{}    // Reset for each port
		ipsWithoutPortOpened = []string{} // Reset for each port

		for ip, openPorts := range openPortsPerIp { // Iterate over ips and check whether this particular port is open
			exists := false
			for _, openPort := range openPorts { // Search data for whether ip has this particular port open
				if openPort == uniqueRelevantPort {
					exists = true
					break
				}
			}
			if exists { // This ip has this port open
				ipsWithPortOpened = append(ipsWithPortOpened, ip)
			} else { // This ip doesn't have this port open
				ipsWithoutPortOpened = append(ipsWithPortOpened, ip)
			}
		}

		if len(ipsWithoutPortOpened) == 0 { // Consistently open port
			openPorts = append(openPorts, uniqueRelevantPort)
		} else { // Inconsistent port
			inconsistencies = append(inconsistencies,
				"Inconsistent open port deteced: Port "+strconv.FormatUint(uint64(uniqueRelevantPort), 10)+
					" is open at "+strings.Join(ipsWithPortOpened, ", ")+
					" but closed at "+strings.Join(ipsWithoutPortOpened, ", "))
		}
	}

	return openPorts, inconsistencies
}
