package dnsScan

import (
	"fmt"
	"io"
	"log/slog"
)

func PrintMailConfigScanResults(results []string, out io.Writer) {

	slog.Debug("dnsScan: Printing mail config scan results started")

	if len(results) > 0 {

		if _, err := fmt.Fprintf(out, "\n\n## Mail security scan results\n"); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}

		for _, message := range results {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("dnsScan: Error writing to output", "error", err)
			}
		}
	}

	slog.Debug("dnsScan: Printing mail config scan results completed")

}
