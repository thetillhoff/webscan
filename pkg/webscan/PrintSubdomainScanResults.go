package webscan

import (
	"fmt"
)

func (engine Engine) PrintSubdomainScanResults() {

	if engine.SubdomainScan && len(engine.subdomains) > 0 {

		fmt.Printf("\n\n--- Subdomain scan results ---\n")

		fmt.Println("Subdomains found in certificate logs (crt.sh):")
		for index, subDomainName := range engine.subdomains {
			fmt.Println(index, subDomainName)
		}
	}

}
