package tlsScan

import (
	"crypto/tls"
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func scanEnabledTlsCiphers(status *status.Status, target types.Target, ip string) []tls.CipherSuite {
	var (
		ciphers = []tls.CipherSuite{}

		enabledTlsCiphersChan chan tls.CipherSuite

		enabledTlsCiphers = []tls.CipherSuite{}
	)

	// Create list of all ciphers
	for _, cipher := range tls.CipherSuites() { // Add the secure ciphers
		ciphers = append(ciphers, *cipher)
	}
	for _, cipher := range tls.InsecureCipherSuites() { // Add the insecure ciphers
		ciphers = append(ciphers, *cipher)
	}

	status.SpinningXOfInit(len(ciphers), "Scanning enabled tls ciphers...")

	enabledTlsCiphersChan = make(chan tls.CipherSuite, len(ciphers))

	slog.Debug("tlsScan: Get enabled tls ciphers started", "len(ciphers)", len(ciphers))

	status.SpinningXOfInit(len(ciphers), "Scanning tls ciphers...")

	for _, cipher := range ciphers { // For each cipher
		wg.Add(1)                                                         // Wait for one more goroutine to finish
		go checkCipher(status, target, ip, cipher, enabledTlsCiphersChan) // Start goroutine that checks if tlsVersion and cipher combination are enabled
	}

	wg.Wait()                    // Wait until all goroutines are finished
	close(enabledTlsCiphersChan) // Make sure channel is closed when goroutines are finished
	status.SpinningXOfComplete("Scan of enabled tls ciphers completed.")

	for enabledTlsCipher := range enabledTlsCiphersChan {
		enabledTlsCiphers = append(enabledTlsCiphers, enabledTlsCipher)
	}

	slog.Debug("tlsScan: Get available tls ciphers completed")

	return enabledTlsCiphers
}
