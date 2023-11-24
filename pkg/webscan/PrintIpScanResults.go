package webscan

import (
	"fmt"
	"strings"
)

func (engine Engine) PrintIpScanResults() {

	fmt.Printf("\n\n--- IP scan results ---\n")

	if engine.IpScan && len(engine.ipIsBlacklistedAt) > 0 {
		for ip, blacklists := range engine.ipIsBlacklistedAt {
			fmt.Println(ip, "is blacklisted at", strings.Join(blacklists, ", "))
		}
	} else {
		if len(engine.dnsScanEngine.ARecords) == 1 {
			fmt.Println("IPv4 address is not blacklisted.")
		} else if len(engine.dnsScanEngine.ARecords) > 1 {
			fmt.Println("IPv4 addresses are not blacklisted.")
		}

		if len(engine.dnsScanEngine.AAAARecords) == 1 {
			fmt.Println("IPv6 address is not blacklisted.")
		} else if len(engine.dnsScanEngine.AAAARecords) > 1 {
			fmt.Println("IPv6 addresses are not blacklisted.")
		}
	}

	if engine.IpScan && len(engine.ipOwners) > 0 {
		for _, message := range engine.ipOwners {
			fmt.Println(message)
		}
	}
}
