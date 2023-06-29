package webscan

import (
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine Engine) PrintTlsScanResults() {

	if engine.TlsScan && engine.tlsCiphers != nil {
		tlsScan.PrintRecommendations(engine.tlsCiphers)
	}

}
