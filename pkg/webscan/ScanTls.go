package webscan

import (
	"fmt"

	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine Engine) ScanTls(inputUrl string) (Engine, error) {

	fmt.Println("Scanning TLS...")

	if engine.isAvailableViaHttps {
		engine.tlsResult = tlsScan.ValidateTlsCertificate(inputUrl)
		engine.tlsCiphers = tlsScan.GetAvailableTlsCiphers(inputUrl)
	} // else there is no TLS to be scanned

	return engine, nil
}
