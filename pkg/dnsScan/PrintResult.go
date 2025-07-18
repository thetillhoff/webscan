package dnsScan

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, out io.Writer) {

	slog.Debug("dnsScan: Printing result started")

	if _, err := fmt.Fprintf(out, "## DNS scan results\n\n"); err != nil {
		slog.Debug("dnsScan: Error writing to output", "error", err)
	}

	if len(result.DomainOwners) > 0 {
		if _, err := fmt.Fprintf(out, "Domain Registrar: %s\n", strings.Join(result.DomainOwners, ", ")); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}
	} else {
		slog.Debug("dnsScan: No domain owners found")
	}

	if len(result.NameserverOwners) > 0 {
		if _, err := fmt.Fprintf(out, "Nameserver Owner: %s\n", strings.Join(result.NameserverOwners, ", ")); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}
	} else {
		slog.Debug("dnsScan: No nameserver owners found")
	}

	if len(result.DomainIsBlacklistedAt) > 0 {
		if _, err := fmt.Fprintf(out, "Domain is blacklisted at:\n"); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}
		for _, blacklist := range result.DomainIsBlacklistedAt {
			if _, err := fmt.Fprintf(out, "%s\n", blacklist); err != nil {
				slog.Debug("dnsScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("dnsScan: No blacklist entries found")
	}

	if _, err := fmt.Fprintf(out, "DNS records:\n"); err != nil {
		slog.Debug("dnsScan: Error writing to output", "error", err)
	}
	result.PrintAllDnsRecords(out)

	// Domain Accessibility
	if len(result.OpinionatedHints) > 0 {
		if _, err := fmt.Fprintf(out, "\n"); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}
		for _, hint := range result.OpinionatedHints {
			if _, err := fmt.Fprintf(out, "%s\n", hint); err != nil {
				slog.Debug("dnsScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("dnsScan: No opinionated hints found")
	}

	slog.Debug("dnsScan: Printing result completed")
}
