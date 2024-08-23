package dnsScan

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
)

func ScanMailConfig(resolver *net.Resolver, inputUrl string, txtRecords []string, dkimSelector string) ([]string, error) {
	var (
		results = []string{}
	)

	slog.Debug("dnsScan: Scanning mail config started")

	fmt.Println("Scanning mail config...")

	// if engine.SubdomainScan {
	if dkimSelector != "" {
		results = CheckMailConfig(inputUrl, resolver, txtRecords, dkimSelector)
	} else {
		return results, errors.New("DKIM selector required")
	}
	// }

	slog.Debug("dnsScan: Scanning mail config completed")

	return results, nil
}
