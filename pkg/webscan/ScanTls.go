package webscan

import (
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine Engine) ScanTls() (Engine, error) {

	if engine.isAvailableViaHttps {
		engine.tlsResult = tlsScan.ValidateTlsCertificate(engine.url)
		engine.tlsCiphers = tlsScan.GetAvailableTlsCiphers(engine.url)
	} // else there is no TLS to be scanned

	return engine, nil
}
