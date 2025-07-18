package httpHeaderScan

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, schemaName string, out io.Writer) {
	var (
		messages = []string{}
	)

	slog.Debug("httpHeaderScan: Printing result started")

	messages = append(messages, result.httpHeaderRecommendations...)

	if len(result.httpCookieRecommendations)+len(result.httpOtherCookieRecommendations) > 0 { // If any recommendations for cookies exist
		messages = append(messages, "\nCookies:\n") // Add empty line and subheading for better readability
	}

	for cookieName, recommendations := range result.httpCookieRecommendations {
		for _, recommendation := range recommendations {
			messages = append(messages, fmt.Sprintf("Cookie '%s' %s", cookieName, recommendation))
		}
	}

	messages = append(messages, result.httpOtherCookieRecommendations...)

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## %s header scan results\n\n", strings.ToUpper(schemaName)); err != nil {
			slog.Debug("httpHeaderScan: Error writing to output", "error", err)
		}

		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("httpHeaderScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("httpHeaderScan: No information found")
	}

	slog.Debug("httpHeaderScan: Printing result completed")

}
