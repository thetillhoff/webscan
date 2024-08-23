package subDomainScan

import (
	"log/slog"
	"strings"

	"github.com/thetillhoff/webscan/pkg/status"
)

func Scan(status *status.Status, inputUrl string, certNames map[string]struct{}) Result {
	var (
		err    error
		result = Result{
			subdomainsFromTlsScan: map[string]struct{}{},
			subdomainsFromCrtSh:   map[string]struct{}{},
		}
	)

	slog.Debug("subDomainScan: Scan started")

	status.SpinningUpdate("Scanning subdomains...")

	for subdomain := range certNames {

		if strings.Contains(subdomain, "=") { // If "subdomain" contains strings like `CN=X,O=Y,C=Z`
			continue // Skip this entry
		}

		subdomain = strings.TrimPrefix(subdomain, "*.") // Remove wildcards, as they are invalid dns names, but might contain valid subdomains
		if subdomain != inputUrl {                      // Remove names that equal the input domain
			if _, ok := result.subdomainsFromTlsScan[subdomain]; !ok { // Skip duplicates
				result.subdomainsFromTlsScan[subdomain] = struct{}{} // Add unique entries
			}
		}
	}

	result.subdomainsFromCrtSh, err = CheckCertLogs(inputUrl)
	if err != nil {
		slog.Warn("could not retrieve subdomains from crt.sh")
	}

	status.SpinningComplete("Scan of subdomains complete.")

	slog.Debug("subDomainScan: Scan completed")

	return result
}
