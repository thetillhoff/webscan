package webscan

import "fmt"

func (engine Engine) PrintPortScanResults() {

	if engine.DetailedPortScan {

		fmt.Printf("\n\n--- TCP port scan results ---\n")

		if len(engine.openPorts) > 0 {
			fmt.Println("Relevant open ports:")
			for _, relevantOpenPort := range engine.openPorts {
				fmt.Println("-", relevantOpenPort)
			}
		} else {
			fmt.Println("No relevant open ports found.")
		}
	}

	if len(engine.openPortInconsistencies) > 0 {
		for _, portInconsistency := range engine.openPortInconsistencies {
			fmt.Println(portInconsistency)
		}
	}

}
