package webscan

import "fmt"

func (engine Engine) PrintPortScanResults() {

	if engine.DetailedPortScan {
		if len(engine.portScanOpenPorts) > 0 {
			fmt.Println()
			fmt.Println("Relevant open ports:")
			for _, relevantOpenPort := range engine.portScanOpenPorts {
				fmt.Println("-", relevantOpenPort)
			}
		} else {
			fmt.Println()
			fmt.Println("No relevant open ports found.")
		}
	}

}
