package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/thetillhoff/webscan/pkg/status"
)

func checkCipher(status *status.Status, url string, tlsCipher tls.CipherSuite, allowedCiphers chan<- tls.CipherSuite) {

	defer wg.Done()
	defer status.SpinningXOfUpdate()

	slog.Debug("tlsScan: Checking if cipher is available started", "cipher", tlsCipher.Name)

	_, err := tls.Dial("tcp", url+":443", &tls.Config{
		MinVersion:       tls.VersionTLS10,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ServerName:       url,
		CipherSuites: []uint16{
			tlsCipher.ID,
		},
	})

	// TODO try each TLS version
	// TODO try each cipher, warn on insecure, at least one secure

	if !os.IsTimeout(err) && err == nil { // If no timeout error occurred and there was no other error
		allowedCiphers <- tlsCipher
	}

	slog.Debug("tlsScan: Checking if cipher is available completed", "cipher", tlsCipher.Name)
}
