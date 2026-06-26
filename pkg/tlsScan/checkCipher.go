package tlsScan

import (
	"crypto/tls"
	"log/slog"
	"os"
	"sync"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

func checkCipher(status *status.Status, target types.Target, ip string, tlsCipher tls.CipherSuite, allowedCiphers chan<- tls.CipherSuite, wg *sync.WaitGroup) {
	var targetEndpoint = ip + ":" + target.Port()

	defer wg.Done()
	defer status.SpinningXOfUpdate()

	slog.Debug("tlsScan: Checking if cipher is available started", "targetEndpoint", targetEndpoint, "cipher", tlsCipher.Name)

	conn, err := tls.DialWithDialer(dialer, "tcp", targetEndpoint, &tls.Config{
		MinVersion:       tls.VersionTLS10,
		MaxVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ServerName:       target.Hostname(),
		CipherSuites: []uint16{
			tlsCipher.ID,
		},
	})

	if !os.IsTimeout(err) && err == nil {
		if closeErr := conn.Close(); closeErr != nil {
			slog.Debug("tlsScan: Error closing connection", "error", closeErr)
		}
		allowedCiphers <- tlsCipher
	}

	slog.Debug("tlsScan: Checking if cipher is available completed", "cipher", tlsCipher.Name, "available", err == nil)
}
