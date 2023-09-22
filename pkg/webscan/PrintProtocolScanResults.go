package webscan

import (
	"fmt"
)

func (engine Engine) PrintProtocolScanResults() {

	if len(engine.protocolRecommendations) > 0 {
		fmt.Println()
		for _, message := range engine.protocolRecommendations {
			fmt.Println(message)
		}
	}
}
