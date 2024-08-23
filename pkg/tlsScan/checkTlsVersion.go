package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

func checkTlsVersion(status *status.Status, url string, tlsVersion uint16, allowedTlsVersions chan<- uint16) {

	defer wg.Done()
	defer status.SpinningXOfUpdate()

	slog.Debug("tlsScan: Checking if tls version is available started", "tlsVersion", tls.VersionName(tlsVersion))

	_, err := tls.Dial("tcp", url+":443", &tls.Config{
		MinVersion:       tlsVersion,
		MaxVersion:       tlsVersion,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ServerName:       url,
	})

	// TODO try each TLS version
	// TODO try each cipher, warn on insecure, at least one secure

	if !os.IsTimeout(err) && err == nil { // If no timeout error occurred and there was no other error
		allowedTlsVersions <- tlsVersion
	}

	slog.Debug("tlsScan: Checking if cipher is available completed", "tlsVersion", tls.VersionName(tlsVersion))
}
