package webscan

import "fmt"

func (engine Engine) PrintPortScanResults() {

	if engine.DetailedPortScan {
		fmt.Println()
		if len(engine.portScanOpenPorts) > 0 {
			fmt.Println("Relevant open ports:")
			for _, relevantOpenPort := range engine.portScanOpenPorts {
				fmt.Println("-", relevantOpenPort)
			}
		} else {
			fmt.Println("No relevant open ports found.")
		}
	}

	if len(engine.portScanInconsistencies) > 0 {
		fmt.Println()
		for _, portInconsistency := range engine.portScanInconsistencies {
			fmt.Println(portInconsistency)
		}
	}

}
