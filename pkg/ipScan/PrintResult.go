package ipScan

import (
	"fmt"
	"log/slog"
	"strings"
)

func PrintResult(result Result, aRecords []string, aaaaRecords []string) {

	slog.Debug("ipScan: Printing result started")

	fmt.Printf("\n\n## IP scan results\n\n")

	if len(result.IpIsBlacklistedAt) > 0 {
		for ip, blacklists := range result.IpIsBlacklistedAt {
			fmt.Println(ip, "is blacklisted at", strings.Join(blacklists, ", ")+"!")
		}
	} else {
		// TODO the following messages should only be printed on DBG/INF level
		// TODO and could be replaced with "no blacklist entries for ip <ip>" or something similar
		if len(aRecords) == 1 {
			fmt.Println("No blocklist entry for IPv4 address found")
		} else if len(aRecords) > 1 {
			fmt.Println("No blocklist entry for IPv4 addresses found.")
		}

		if len(aaaaRecords) == 1 {
			fmt.Println("No blocklist entry for IPv6 address found.")
		} else if len(aaaaRecords) > 1 {
			fmt.Println("No blocklist entry for IPv6 addresses found.")
		}
	}

	if len(result.IpOwners) > 0 {
		for _, message := range result.IpOwners {
			fmt.Println(message)
		}
	}

	slog.Debug("ipScan: Printing result completed")

}
