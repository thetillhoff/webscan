package webscan

import (
	"fmt"
	"strings"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
)

func (engine Engine) PrintDnsScanResults(inputUrl string) {
	if engine.DetailedDnsScan {
		dnsServer := engine.DnsScanEngine.GetCustomDns()
		if dnsServer != "" {
			fmt.Println("Using custom dns server:", dnsServer)
		} else {
			fmt.Println("Using system dns server")
		}

		fmt.Println()
		if len(engine.DnsScanEngine.DomainOwners) == 0 {
			fmt.Println("Could not retrieve Domain Owner (country TLDs are not supported yet by RDAP)")
		} else {
			fmt.Println("Domain Registrar: ", strings.Join(engine.DnsScanEngine.DomainOwners, ", "))
		}

		fmt.Println()
		fmt.Println("DNS records:")
		engine.DnsScanEngine.PrintAllDnsRecords()

		if engine.Opinionated {
			// Domain Accessibility
			fmt.Println()
			for _, hint := range dnsScan.GetDomainAccessibilityHints(inputUrl) {
				fmt.Println(hint)
			}
		}
	}
}
