package webscan

import "fmt"

func (engine Engine) PrintIpScanResults() {

	if engine.IpScan && len(engine.ipScanResult) > 0 {
		fmt.Println()
		for _, message := range engine.ipScanResult {
			fmt.Println(message)
		}
	}

	if engine.IpScan && len(engine.ipScanOwners) > 0 {
		fmt.Println()
		for _, message := range engine.ipScanOwners {
			fmt.Println(message)
		}
	}

}
