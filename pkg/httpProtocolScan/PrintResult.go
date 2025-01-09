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

		if len(result.httpVersions) > 0 {
			fmt.Printf("\nThe following protocols are available for HTTP:\n")
			for _, version := range result.httpVersions {
				if version != "" {
					fmt.Println(version)
				}
			}
		}

		if len(result.httpsVersions) > 0 {
			fmt.Printf("\nThe following protocols are available for HTTPS:\n")
			for _, version := range result.httpsVersions {
				if version != "" {
					fmt.Println(version)
				}
			}
		}
	}

	slog.Debug("httpProtocolScan: Printing result completed")

}
