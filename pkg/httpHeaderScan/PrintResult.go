package httpHeaderScan

import (
	"fmt"
	"log/slog"
	"strings"
)

func PrintResult(result Result, schemaName string) {

	slog.Debug("httpHeaderScan: Printing result started")

	if (len(result.httpHeaderRecommendations) +
		len(result.httpCookieRecommendations) +
		len(result.httpOtherCookieRecommendations)) > 0 {

		fmt.Printf("\n\n## %s header scan results\n\n", strings.ToUpper(schemaName))

		for _, message := range result.httpHeaderRecommendations {
			fmt.Println(message)
		}

		if len(result.httpCookieRecommendations)+len(result.httpOtherCookieRecommendations) > 0 { // If any recommendations for cookies exist
			fmt.Printf("\nCookies:\n") // Add empty line and subheading for better readability
		}

		for cookieName, recommendations := range result.httpCookieRecommendations {
			for _, recommendation := range recommendations {
				fmt.Println("Cookie '" + cookieName + "' " + recommendation)
			}
		}

		for _, message := range result.httpOtherCookieRecommendations {
			fmt.Println(message)
		}
	}

	slog.Debug("httpHeaderScan: Printing result completed")

}
