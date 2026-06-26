package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"sync"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

func scanEnabledTLSVersions(status *status.Status, target types.Target, ip string) []uint16 {
	var (
		wg sync.WaitGroup

		tlsVersions = []uint16{
			tls.VersionTLS13,
			tls.VersionTLS12,
			tls.VersionTLS11,
			tls.VersionTLS10,
		}

		enabledTlsVersionsChan = make(chan uint16, len(tlsVersions))

		enabledTlsVersions = []uint16{}
	)

	slog.Debug("tlsScan: Scanning enabled tls versions started", "len(versions)", len(tlsVersions))

	status.SpinningXOfInit(len(tlsVersions), "Scanning enabled tls versions...")

	for _, tlsVersion := range tlsVersions {
		wg.Add(1)
		go checkTLSVersion(status, target, ip, tlsVersion, enabledTlsVersionsChan, &wg)
	}

	wg.Wait()
	close(enabledTlsVersionsChan)
	status.SpinningXOfComplete("Scan of enabled tls versions completed.")

	for enabledTlsVersion := range enabledTlsVersionsChan {
		enabledTlsVersions = append(enabledTlsVersions, enabledTlsVersion)
	}

	slog.Debug("tlsScan: Scanning enabled tls versions completed")

	return enabledTlsVersions
}
