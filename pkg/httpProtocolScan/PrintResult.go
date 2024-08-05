package httpProtocolScan

import (
	"fmt"
	"log/slog"
)

func PrintResult(result Result) {

	slog.Debug("httpProtocolScan: Printing result started")

	if len(result.recommendations) > 0 {

		fmt.Printf("\n\n## HTTP protocol scan results\n\n")

		for _, message := range result.recommendations {
			fmt.Println(message)
		}
	}

	slog.Debug("httpProtocolScan: Printing result completed")

}
