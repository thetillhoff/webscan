package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"sync"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func scanEnabledTLSCiphers(status *status.Status, target types.Target, ip string) []tls.CipherSuite {
	var (
		wg sync.WaitGroup

		ciphers = []tls.CipherSuite{}

		enabledTlsCiphersChan chan tls.CipherSuite

		enabledTlsCiphers = []tls.CipherSuite{}
	)

	// Create list of all ciphers
	for _, cipher := range tls.CipherSuites() {
		ciphers = append(ciphers, *cipher)
	}
	for _, cipher := range tls.InsecureCipherSuites() {
		ciphers = append(ciphers, *cipher)
	}

	enabledTlsCiphersChan = make(chan tls.CipherSuite, len(ciphers))

	slog.Debug("tlsScan: Scanning enabled TLS ciphers started", "count", len(ciphers))

	status.SpinningXOfInit(len(ciphers), "Scanning tls ciphers...")

	for _, cipher := range ciphers {
		wg.Add(1)
		go checkCipher(status, target, ip, cipher, enabledTlsCiphersChan, &wg)
	}

	wg.Wait()
	close(enabledTlsCiphersChan)
	status.SpinningXOfComplete("Scan of enabled tls ciphers completed.")

	for enabledTlsCipher := range enabledTlsCiphersChan {
		enabledTlsCiphers = append(enabledTlsCiphers, enabledTlsCipher)
	}

	slog.Debug("tlsScan: Scanning enabled TLS ciphers completed")

	return enabledTlsCiphers
}
