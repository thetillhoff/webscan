package httpHeaderScan

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, schemaName string, out io.Writer) {
	hasEntries := len(result.httpHeaderEntries) > 0
	hasCookies := len(result.httpCookieRecommendations)+len(result.httpOtherCookieRecommendations) > 0

	slog.Debug("httpHeaderScan: Printing result started")

	if !hasEntries && !hasCookies {
		slog.Debug("httpHeaderScan: No information found")
		return
	}

	fmt.Fprintf(out, "\n## %s header scan results\n\n", strings.ToUpper(schemaName))

	if hasEntries {
		fmt.Fprintf(out, "Headers:\n")
		for _, entry := range result.httpHeaderEntries {
			if strings.Contains(entry.Value, "\n") {
				fmt.Fprintf(out, "- %s:\n", entry.Name)
				for _, line := range strings.Split(entry.Value, "\n") {
					fmt.Fprintf(out, "    %s\n", line)
				}
			} else if entry.Value == "" {
				fmt.Fprintf(out, "- %s: (not set)\n", entry.Name)
			} else {
				fmt.Fprintf(out, "- %s: %s\n", entry.Name, entry.Value)
			}
			if entry.Recommendation != "" {
				fmt.Fprintf(out, "  → %s\n", entry.Recommendation)
			}
		}
	}

	if hasCookies {
		fmt.Fprintf(out, "\nCookies:\n")
		for cookieName, recommendations := range result.httpCookieRecommendations {
			for _, recommendation := range recommendations {
				fmt.Fprintf(out, "- Cookie '%s' %s\n", cookieName, recommendation)
			}
		}
		for _, rec := range result.httpOtherCookieRecommendations {
			fmt.Fprintf(out, "- %s\n", rec)
		}
	}

	slog.Debug("httpHeaderScan: Printing result completed")
}
