package ipScan

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, aRecords []string, aaaaRecords []string, out io.Writer) {
	var (
		messages = []string{}
	)

	slog.Debug("ipScan: Printing result started")

	if len(result.IpIsBlacklistedAt) > 0 {
		for ip, blacklists := range result.IpIsBlacklistedAt {
			messages = append(messages, fmt.Sprintf("%s is blacklisted at %s!", ip, strings.Join(blacklists, ", ")))
		}
	}

	if len(result.IpOwners) > 0 {
		messages = append(messages, result.IpOwners...)
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## IP scan results\n\n"); err != nil {
			slog.Debug("ipScan: Error writing to output", "error", err)
		}
		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("ipScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("ipScan: No information found")
	}

	slog.Debug("ipScan: Printing result completed")

}
