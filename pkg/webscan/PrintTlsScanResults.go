package webscan

import (
	"fmt"

	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

func (engine Engine) PrintTlsScanResults() {

	if engine.TlsScan {

		fmt.Printf("\n\n--- TLS scan results ---\n")

		if engine.tlsResult != nil {
			fmt.Println(engine.tlsResult)
		}

		if engine.tlsCiphers != nil {
			for _, message := range tlsScan.GetRecommendations(engine.tlsCiphers) {
				fmt.Println(message)
			}
		}
	}
}
