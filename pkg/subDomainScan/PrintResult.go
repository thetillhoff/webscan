package subDomainScan

import (
	"fmt"
	"io"
	"log/slog"
)

func PrintResult(result Result, out io.Writer) {
	var (
		messages = []string{}

		maxLength = 0
	)

	slog.Debug("subDomainScan: Printing result started")

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
		messages = append(messages, fmt.Sprintf("- %*s (from tls certificate)", maxLength, subDomainName))
	}

	for subDomainName := range result.subdomainsFromCrtSh {
		messages = append(messages, fmt.Sprintf("- %*s (from crt.sh)", maxLength, subDomainName))
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n\n## Subdomain scan results\n\n"); err != nil {
			slog.Debug("subDomainScan: Error writing to output", "error", err)
		}
		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("subDomainScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("subDomainScan: No information found")
	}

	slog.Debug("subDomainScan: Printing result completed")

}
