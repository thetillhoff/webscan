package webscan

import (
	"fmt"
	"log"

	"github.com/thetillhoff/webscan/pkg/ipScan"
)

func (engine Engine) ScanIps() (Engine, error) {
	var (
		err error

		response string
	)

	if (len(engine.dnsScanEngine.ARecords) + len(engine.dnsScanEngine.AAAARecords)) > 1 { // If there is more than one IP
		fmt.Println("Scanning IPs...") // Plural
	} else {
		fmt.Println("Scanning IP...") // Singular
	}

	for _, aRecord := range engine.dnsScanEngine.ARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aRecord)
		if err != nil {
			log.Fatalln(err)
		}
		engine.ipScanOwners = append(engine.ipScanOwners, "According to RDAP information, IP "+aRecord+" is registered at "+response)
	}

	for _, aaaaRecord := range engine.dnsScanEngine.AAAARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aaaaRecord)
		if err != nil {
			log.Fatalln(err)
		}
		engine.ipScanOwners = append(engine.ipScanOwners, "According to RDAP information, IP "+aaaaRecord+" is registered at "+response)
	}

	return engine, nil
}
