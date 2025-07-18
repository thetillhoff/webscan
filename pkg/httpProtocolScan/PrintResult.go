package httpProtocolScan

import (
	"fmt"
	"io"
	"log/slog"
)

func PrintResult(result Result, out io.Writer) {
	var (
		messages = []string{}
	)

	slog.Debug("httpProtocolScan: Printing result started")

	if len(result.recommendations) > 0 {
		messages = append(messages, "\nThe following recommendations were found:")
		for _, recommendation := range result.recommendations {
			messages = append(messages, fmt.Sprintf("- %s", recommendation))
		}
	}

	if len(result.httpVersions) > 0 {
		messages = append(messages, "\nThe following protocols are available for HTTP:")
		for _, version := range result.httpVersions {
			if version != "" {
				messages = append(messages, fmt.Sprintf("- %s", version))
			}
		}
	}

	if len(result.httpsVersions) > 0 {
		messages = append(messages, "\nThe following protocols are available for HTTPS:")
		for _, version := range result.httpsVersions {
			if version != "" {
				messages = append(messages, fmt.Sprintf("- %s", version))
			}
		}
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## HTTP protocol scan results\n"); err != nil {
			slog.Debug("httpProtocolScan: Error writing to output", "error", err)
		}
		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("httpProtocolScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("httpProtocolScan: No information found")
	}

	slog.Debug("httpProtocolScan: Printing result completed")

}
