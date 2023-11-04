package webscan

import "fmt"

func (engine Engine) PrintHttpHeaderScanResults() {

	if len(engine.httpHeaderRecommendations) > 0 {

		fmt.Printf("\n\n--- HTTP header scan results ---\n")

		for _, message := range engine.httpHeaderRecommendations {
			fmt.Println(message)
		}

		if len(engine.httpCookieRecommendations)+len(engine.httpOtherCookieRecommendations) > 0 { // If any recommendations for cookies exist
			fmt.Println() // Add empty line for better readability
		}

		for cookieName, recommendations := range engine.httpCookieRecommendations {
			for _, recommendation := range recommendations {
				fmt.Println("Cookie '" + cookieName + "' " + recommendation)
			}
		}

		for _, message := range engine.httpOtherCookieRecommendations {
			fmt.Println(message)
		}
	}

}
