package dnsScan

import (
	"fmt"
	"log/slog"
)

func PrintMailConfigScanResults(results []string) {

	slog.Debug("dnsScan: Printing mail config scan results started")

	if len(results) > 0 {

		fmt.Printf("\n\n## Mail security scan results\n")

		for _, message := range results {
			fmt.Println(message)
		}
	}

	slog.Debug("dnsScan: Printing mail config scan results completed")

}
