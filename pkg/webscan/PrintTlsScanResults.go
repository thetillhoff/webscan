package webscan

import (
	"fmt"

	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine Engine) PrintTlsScanResults() {

	if engine.tlsResult != nil {
		fmt.Println()
		fmt.Println(engine.tlsResult)
	}

	if engine.TlsScan && engine.tlsCiphers != nil {
		tlsScan.PrintRecommendations(engine.tlsCiphers)
	}

}
