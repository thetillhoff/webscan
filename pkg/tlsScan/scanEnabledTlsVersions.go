package tlsScan

import (
	"crypto/tls"
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func scanEnabledTlsVersions(status *status.Status, target types.Target, ip string) []uint16 {
	var (
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

	for _, tlsVersion := range tlsVersions { // For each cipher
		wg.Add(1)                                                                  // Wait for one more goroutine to finish
		go checkTlsVersion(status, target, ip, tlsVersion, enabledTlsVersionsChan) // Start goroutine that checks if tlsVersion is enabled
	}

	wg.Wait()                     // Wait until all goroutines are finished
	close(enabledTlsVersionsChan) // Make sure channel is closed when goroutines are finished
	status.SpinningXOfComplete("Scan of enabled tls versions completed.")

	for enabledTlsVersion := range enabledTlsVersionsChan {
		enabledTlsVersions = append(enabledTlsVersions, enabledTlsVersion)
	}

	slog.Debug("tlsScan: Get available tls ciphers completed")

	return enabledTlsVersions
}
