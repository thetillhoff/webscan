package webscan

import "fmt"

func (engine Engine) PrintMailConfigScanResults() {

	if len(engine.mailConfigRecommendations) > 0 {
		fmt.Println()
		fmt.Println("\nDNS mail security findings:")
		for _, message := range engine.mailConfigRecommendations {
			fmt.Println(message)
		}
	}

}
