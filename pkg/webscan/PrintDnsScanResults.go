package webscan

import (
	"fmt"
	"strings"
)

func (engine Engine) PrintDnsScanEngines() {
	if engine.DetailedDnsScan {

		fmt.Println()
		if len(engine.dnsScanEngine.DomainOwners) == 0 {
			fmt.Println("Could not retrieve Domain Owner (country TLDs are not supported yet by RDAP)")
		} else {
			fmt.Println("Domain Registrar: ", strings.Join(engine.dnsScanEngine.DomainOwners, ", "))
		}

		fmt.Println()
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
