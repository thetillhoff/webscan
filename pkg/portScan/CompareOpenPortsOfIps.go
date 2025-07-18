package portScan

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

// Checks whether the open ports are the same of all ip addresses
// Returns list of ports that where open on all ip addresses,
// and a list of inconsistencies
func CompareOpenPortsOfIps(openPortsPerIp map[string][]uint16) ([]uint16, []string) {
	var (
		uniqueRelevantPorts       = map[uint16]struct{}{}
		sortedUniqueRelevantPorts []uint16
		ipsWithPortOpened         []string
		ipsWithoutPortOpened      []string

		openPorts       = []uint16{}
		inconsistencies = []string{}
	)

	slog.Debug("portScan: Comparison of open ports of ips started")

	// Create map where each uint16 value of all ips is an empty struct -> result is a list of unique relevant ports
	for _, openPorts := range openPortsPerIp { // Iterate over ips
		for _, openPort := range openPorts { // Iterate over openPorts of ips
			uniqueRelevantPorts[openPort] = struct{}{} // Add openPort to list of unique relevant ports
		}
	}

	// Sorting unique open ports
	sortedUniqueRelevantPorts = make([]uint16, len(uniqueRelevantPorts)) // Allocate slice with correct length
	portIndex := 0
	for key := range uniqueRelevantPorts { // For each port in uniqueRelevantPorts
		sortedUniqueRelevantPorts[portIndex] = key // Add openPort to slice
		portIndex++                                // Increase index in slice
	}
	slices.Sort(sortedUniqueRelevantPorts) // Sort slice

	// compare open ports per ip against global list of open ports
	for _, uniqueRelevantPort := range sortedUniqueRelevantPorts { // Iterate over unique relevant ports and compare to actual open ports of each ip
		ipsWithPortOpened = []string{}    // Reset for each iteration (==port)
		ipsWithoutPortOpened = []string{} // Reset for each iteration (==port)

		for ip, openPorts := range openPortsPerIp { // Iterate over ips to get their open ports
			portOpenAtThisIp := false
			for _, openPort := range openPorts { // Search data for whether ip has this particular port open
				if openPort == uniqueRelevantPort { // If ip has this particular port open
					portOpenAtThisIp = true
					break // No need to check other ports
				}
			}
			if portOpenAtThisIp { // If this ip has this particular port open
				ipsWithPortOpened = append(ipsWithPortOpened, ip)
			} else { // If this ip doesn't have this particular port open
				ipsWithoutPortOpened = append(ipsWithoutPortOpened, ip)
			}
		}

		openPorts = append(openPorts, uniqueRelevantPort) // Add open port of any ip address to list of open ports, even if inconsistently open

		if len(ipsWithoutPortOpened) > 0 { // Inconsistent port
			inconsistencies = append(inconsistencies,
				"Inconsistent open port detected: Port "+strconv.FormatUint(uint64(uniqueRelevantPort), 10)+
					" is open at "+strings.Join(ipsWithPortOpened, ", ")+
					" but closed at "+strings.Join(ipsWithoutPortOpened, ", "))
		}
	}

	slog.Debug("portScan: Comparison of open ports of ips completed")

	return openPorts, inconsistencies
}
