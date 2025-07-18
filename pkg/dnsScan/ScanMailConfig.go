package dnsScan

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/miekg/dns"
)

func ScanMailConfig(dnsClient *dns.Client, nameserver string, inputUrl string, txtRecords []string, dkimSelector string, out io.Writer) ([]string, error) {
	var (
		results = []string{}
	)

	slog.Debug("dnsScan: Scanning mail config started")

	if _, err := fmt.Fprintf(out, "Scanning mail config...\n"); err != nil {
		slog.Debug("dnsScan: Error writing to output", "error", err)
	}

	// if engine.SubdomainScan {
	if dkimSelector != "" {
		results = CheckMailConfig(inputUrl, dnsClient, nameserver, txtRecords, dkimSelector)
	} else {
		return results, errors.New("DKIM selector required")
	}
	// }

	slog.Debug("dnsScan: Scanning mail config completed")

	return results, nil
}
