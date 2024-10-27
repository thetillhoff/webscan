package dnsScan

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

func PrintResult(result Result) {

	slog.Debug("dnsScan: Printing result started")

	fmt.Printf("\n\n## DNS scan results\n\n")

	if len(result.DomainOwners) == 0 {
		fmt.Println("Could not retrieve Domain Owner (country TLDs are not supported yet by RDAP)")
	} else {
		slices.Sort(result.DomainOwners)
		fmt.Println("Domain Registrar: ", strings.Join(result.DomainOwners, ", "))
	}

	if len(result.NameserverOwners) == 0 {
		fmt.Println("Could not retrieve Nameserver Owner")
	} else {
		slices.Sort(result.NameserverOwners)
		fmt.Println("Nameserver Owner:", strings.Join(result.NameserverOwners, ", "))
	}

	if len(result.DomainIsBlacklistedAt) > 0 {
		fmt.Println("Domain is blacklisted at:")
		for _, blacklist := range result.DomainIsBlacklistedAt {
			fmt.Println(blacklist)
		}
	} else {
		fmt.Println("Domain is not blacklisted.")
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
