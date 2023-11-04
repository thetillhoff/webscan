package webscan

import "fmt"

func (engine Engine) PrintMailConfigScanResults() {

	if len(engine.mailConfigRecommendations) > 0 {

		fmt.Printf("\n\n--- Mail security scan results ---\n")

		for _, message := range engine.mailConfigRecommendations {
			fmt.Println(message)
		}
	}

}
