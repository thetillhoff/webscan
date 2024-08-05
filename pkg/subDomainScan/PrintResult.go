package subDomainScan

import (
	"fmt"
	"log/slog"
)

func PrintResult(result Result) {

	slog.Debug("subDomainScan: Printing result started")

	slog.Debug("subDomain results", "len(result.subdomainsFromTlsScan)", len(result.subdomainsFromTlsScan), "len(result.subdomainsFromCrtSh)", len(result.subdomainsFromCrtSh))

	if len(result.subdomainsFromTlsScan) > 0 || len(result.subdomainsFromCrtSh) > 0 {

		fmt.Printf("\n\n## Subdomain scan results\n\n")

		maxLength := 0
		for subDomainName := range result.subdomainsFromTlsScan {
			if len(subDomainName) > maxLength {
				maxLength = len(subDomainName)
			}
		}
		for subDomainName := range result.subdomainsFromCrtSh {
			if len(subDomainName) > maxLength {
				maxLength = len(subDomainName)
			}
		}

		for subDomainName := range result.subdomainsFromTlsScan {
			fmt.Printf("- %*s (from tls certificate)\n", maxLength, subDomainName)
		}

		for subDomainName := range result.subdomainsFromCrtSh {
			fmt.Printf("- %*s (from crt.sh)\n", maxLength, subDomainName)
		}
	}

	slog.Debug("subDomainScan: Printing result completed")

}
