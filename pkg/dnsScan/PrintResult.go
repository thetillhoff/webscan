package dnsScan

import (
	"fmt"
	"log/slog"
	"strings"
)

func PrintResult(result Result) {

	slog.Debug("dnsScan: Printing result started")

	fmt.Printf("\n\n## DNS scan results\n\n")

	if len(result.DomainOwners) > 0 {
		fmt.Println("Domain Registrar: ", strings.Join(result.DomainOwners, ", "))
	}

	if len(result.NameserverOwners) > 0 {
		fmt.Println("Nameserver Owner:", strings.Join(result.NameserverOwners, ", "))
	}

	if len(result.DomainIsBlacklistedAt) > 0 {
		fmt.Println("Domain is blacklisted at:")
		for _, blacklist := range result.DomainIsBlacklistedAt {
			fmt.Println(blacklist)
		}
	}

	fmt.Println("DNS records:")
	result.PrintAllDnsRecords()

	// Domain Accessibility
	if len(result.OpinionatedHints) > 0 {
		fmt.Println()
		for _, hint := range result.OpinionatedHints {
			fmt.Println(hint)
		}
	}

	slog.Debug("dnsScan: Printing result completed")
}
