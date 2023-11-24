package webscan

import (
	"fmt"

	"github.com/thetillhoff/webscan/pkg/ipScan"
)

func (engine Engine) ScanIps() (Engine, error) {
	var (
		err error

		response         string
		blacklistMatches []string
	)

	if (len(engine.dnsScanEngine.ARecords) + len(engine.dnsScanEngine.AAAARecords)) > 1 { // If there is more than one IP
		fmt.Println("Scanning IPs...") // Plural
	} else {
		fmt.Println("Scanning IP...") // Singular
	}

	for _, aRecord := range engine.dnsScanEngine.ARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aRecord)
		if err != nil {
			return engine, err
		}
		engine.ipOwners = append(engine.ipOwners, "According to RDAP information, IP "+aRecord+" is registered at "+response)

		blacklistMatches, err = ipScan.IsIpBlacklisted(aRecord, engine.Verbose)
		if err != nil {
			return engine, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			engine.ipIsBlacklistedAt[aRecord] = blacklistMatches
		}
	}

	for _, aaaaRecord := range engine.dnsScanEngine.AAAARecords {
		response, err = ipScan.GetIPOwnerViaRDAP(aaaaRecord)
		if err != nil {
			return engine, err
		}
		engine.ipOwners = append(engine.ipOwners, "According to RDAP information, IP "+aaaaRecord+" is registered at "+response)

		blacklistMatches, err = ipScan.IsIpBlacklisted(aaaaRecord, engine.Verbose)
		if err != nil {
			return engine, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			engine.ipIsBlacklistedAt[aaaaRecord] = blacklistMatches
		}
	}

	return engine, nil
}
