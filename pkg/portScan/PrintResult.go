package portScan

import (
	"fmt"
	"log/slog"
)

func PrintResult(result Result) {

	slog.Debug("portScan: Printing result started")

	fmt.Printf("\n\n## TCP port scan results\n\n")

	if len(result.openPorts) > 0 {
		fmt.Println("Relevant open ports:")
		for _, relevantOpenPort := range result.openPorts {
			fmt.Println("-", relevantOpenPort)
		}
	} else {
		fmt.Println("No relevant open ports found.")
	}

	if len(result.openPortInconsistencies) > 0 {
		for _, portInconsistency := range result.openPortInconsistencies {
			fmt.Println(portInconsistency)
		}
	}

	slog.Debug("portScan: Printing result completed")

}
