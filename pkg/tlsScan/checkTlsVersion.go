package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"os"
	"sync"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func checkTLSVersion(status *status.Status, target types.Target, ip string, tlsVersion uint16, allowedTlsVersions chan<- uint16, wg *sync.WaitGroup) {
	var targetEndpoint = ip + ":" + target.Port()

	defer wg.Done()
	defer status.SpinningXOfUpdate()

	slog.Debug("tlsScan: Checking if tls version is available started", "targetEndpoint", targetEndpoint, "serverName", target.Hostname(), "tlsVersion", tls.VersionName(tlsVersion))

	conn, err := tls.Dial("tcp", targetEndpoint, &tls.Config{
		MinVersion:       tlsVersion,
		MaxVersion:       tlsVersion,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ServerName:       target.Hostname(),
	})

	if !os.IsTimeout(err) && err == nil {
		if closeErr := conn.Close(); closeErr != nil {
			slog.Debug("tlsScan: Error closing connection", "error", closeErr)
		}
		allowedTlsVersions <- tlsVersion
	}

	slog.Debug("tlsScan: Checking if tls version is available completed", "tlsVersion", tls.VersionName(tlsVersion), "available", err == nil)
}
