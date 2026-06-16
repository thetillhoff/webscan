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

	if _, err := fmt.Fprintf(out, "\n## %s header scan results\n\n", strings.ToUpper(schemaName)); err != nil {
		slog.Debug("httpHeaderScan: Error writing to output", "error", err)
	}

	if hasEntries {
		if _, err := fmt.Fprintf(out, "Headers:\n"); err != nil {
			slog.Debug("httpHeaderScan: Error writing to output", "error", err)
		}
		for _, entry := range result.httpHeaderEntries {
			if strings.Contains(entry.Value, "\n") {
				if _, err := fmt.Fprintf(out, "- %s:\n", entry.Name); err != nil {
					slog.Debug("httpHeaderScan: Error writing to output", "error", err)
				}
				for _, line := range strings.Split(entry.Value, "\n") {
					if _, err := fmt.Fprintf(out, "    %s\n", line); err != nil {
						slog.Debug("httpHeaderScan: Error writing to output", "error", err)
					}
				}
			} else if entry.Value == "" {
				if _, err := fmt.Fprintf(out, "- %s: (not set)\n", entry.Name); err != nil {
					slog.Debug("httpHeaderScan: Error writing to output", "error", err)
				}
			} else {
				if _, err := fmt.Fprintf(out, "- %s: %s\n", entry.Name, entry.Value); err != nil {
					slog.Debug("httpHeaderScan: Error writing to output", "error", err)
				}
			}
			if entry.Recommendation != "" {
				if _, err := fmt.Fprintf(out, "  → %s\n", entry.Recommendation); err != nil {
					slog.Debug("httpHeaderScan: Error writing to output", "error", err)
				}
			}
		}
	}

	if hasCookies {
		if _, err := fmt.Fprintf(out, "\nCookies:\n"); err != nil {
			slog.Debug("httpHeaderScan: Error writing to output", "error", err)
		}
		for cookieName, recommendations := range result.httpCookieRecommendations {
			for _, recommendation := range recommendations {
				if _, err := fmt.Fprintf(out, "- Cookie '%s' %s\n", cookieName, recommendation); err != nil {
					slog.Debug("httpHeaderScan: Error writing to output", "error", err)
				}
			}
		}
		for _, rec := range result.httpOtherCookieRecommendations {
			if _, err := fmt.Fprintf(out, "- %s\n", rec); err != nil {
				slog.Debug("httpHeaderScan: Error writing to output", "error", err)
			}
		}
	}

	slog.Debug("httpHeaderScan: Printing result completed")
}
