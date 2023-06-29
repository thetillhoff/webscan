package webscan

import "fmt"

func (engine Engine) PrintHttpHeaderScanResults() {

	if len(engine.httpHeaderRecommendations) > 0 {
		fmt.Println()
		for _, message := range engine.httpHeaderRecommendations {
			fmt.Println(message)
		}
	}

}
