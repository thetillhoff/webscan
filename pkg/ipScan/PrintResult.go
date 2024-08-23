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
	}

	if len(result.IpOwners) > 0 {
		for _, message := range result.IpOwners {
			fmt.Println(message)
		}
	}

	slog.Debug("ipScan: Printing result completed")

}
