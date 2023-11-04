package webscan

import (
	"fmt"
	"strings"
)

func (engine Engine) PrintDnsScanEngines() {
	if engine.DetailedDnsScan {

		fmt.Printf("\n\n--- DNS scan results ---\n")

		if len(engine.dnsScanEngine.DomainOwners) == 0 {
			fmt.Println("Could not retrieve Domain Owner (country TLDs are not supported yet by RDAP)")
		} else {
			fmt.Println("Domain Registrar: ", strings.Join(engine.dnsScanEngine.DomainOwners, ", "))
		}

		if len(engine.dnsScanEngine.DomainIsBlacklistedAt) > 0 {
			fmt.Println("Domain is blacklisted at:")
			for _, blacklist := range engine.dnsScanEngine.DomainIsBlacklistedAt {
				fmt.Println(blacklist)
			}
		} else {
			fmt.Println("Domain is not blacklisted.")
		}

		fmt.Println("DNS records:")
		engine.dnsScanEngine.PrintAllDnsRecords()

		// Domain Accessibility
		if len(engine.dnsScanEngine.OpinionatedHints) > 0 {
			fmt.Println()
			for _, hint := range engine.dnsScanEngine.OpinionatedHints {
				fmt.Println(hint)
			}
		}
	}
}
