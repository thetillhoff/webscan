package portScan

import (
	"fmt"
	"io"
	"log/slog"
)

func PrintResult(result Result, out io.Writer) {
	var (
		messages = []string{}
	)

	slog.Debug("portScan: Printing result started")

	if len(result.openPorts) > 0 {
		messages = append(messages, "Relevant open ports:")
		for _, relevantOpenPort := range result.openPorts {
			messages = append(messages, fmt.Sprintf("- %d", relevantOpenPort))
		}
	} else {
		messages = append(messages, "No relevant open ports found.")
	}

	if len(result.openPortInconsistencies) > 0 {
		messages = append(messages, "Open port inconsistencies:")
		messages = append(messages, result.openPortInconsistencies...)
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## TCP port scan results\n\n"); err != nil {
			slog.Debug("portScan: Error writing to output", "error", err)
		}
		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("portScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("portScan: No information found")
	}

	slog.Debug("portScan: Printing result completed")

}
