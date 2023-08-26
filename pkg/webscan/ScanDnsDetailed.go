package webscan

import (
	"fmt"
)

func (engine Engine) ScanDnsDetailed(inputUrl string) (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning DNS (detailed) of", inputUrl, "...")

	if engine.Opinionated {
		// Domain Accessibility
		engine.dnsScanEngine = engine.dnsScanEngine.CheckIpVersionCompatibility(engine.dnsScanEngine.ARecords, engine.dnsScanEngine.AAAARecords)
		engine.dnsScanEngine = engine.dnsScanEngine.GetDomainAccessibilityHints(inputUrl)
	}

	engine.dnsScanEngine, _ = engine.dnsScanEngine.GetDomainOwnerViaRDAP(inputUrl)
	// // `err` is ignored here, as it's okay that it can't be retrieved. It's not a critical error, but an error nonetheless
	// if err != nil {
	// return engine, err
	// }

	engine.dnsScanEngine, err = engine.dnsScanEngine.GetAllRecords(inputUrl)
	if err != nil {
		return engine, err
	}

	if len(engine.dnsScanEngine.ARecords) == 0 && len(engine.dnsScanEngine.AAAARecords) == 0 && engine.FollowRedirects { // If neither A nor AAAA records exist & redirects should be followed
		engine.dnsScanEngine, err = engine.dnsScanEngine.GetCNAMERecord(inputUrl) // Retrieve CNAME record if exists
		if err != nil {
			return engine, err
		}
		if engine.dnsScanEngine.CNAMERecord != "" { // If CNAME record exists
			if engine.Verbose {
				fmt.Println("No A or AAAA records for", inputUrl, ". Following CNAME...")
			}
			engine, err = engine.ScanDnsSimple(engine.dnsScanEngine.CNAMERecord) // Follow CNAME recursively, but only checking A, AAAA, and CNAMEs
			if err != nil {
				return engine, err
			}
		}
	}

	return engine, nil
}
