package webscan

import (
	"log"

	"github.com/thetillhoff/webscan/pkg/ipScan"
)

func (engine Engine) ScanIps() (Engine, error) {
	var (
		err error

		response string
	)

	response = ipScan.CheckIpVersionCompatibility(engine.DnsScanEngine.ARecords, engine.DnsScanEngine.AAAARecords)

	if response != "" {
		engine.ipScanResult = append(engine.ipScanResult, response)
	}

	for _, aRecord := range engine.DnsScanEngine.ARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aRecord)
		if err != nil {
			log.Fatalln(err)
		}
		engine.ipScanOwners = append(engine.ipScanOwners, "According to RDAP information, IP "+aRecord+" is registered at "+response)
	}

	for _, aaaaRecord := range engine.DnsScanEngine.AAAARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aaaaRecord)
		if err != nil {
			log.Fatalln(err)
		}
		engine.ipScanOwners = append(engine.ipScanOwners, "According to RDAP information, IP "+aaaaRecord+" is registered at "+response)
	}

	return engine, nil
}
