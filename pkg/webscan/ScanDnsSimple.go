package webscan

import "fmt"

func (engine Engine) ScanDnsSimple(inputUrl string) (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning DNS (simple) of", inputUrl, "...")

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetARecords(inputUrl)
	if err != nil {
		return engine, err
	}

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetAAAARecords(inputUrl)
	if err != nil {
		return engine, err
	}

	if len(engine.dnsScanResult.ARecords) == 0 && len(engine.dnsScanResult.AAAARecords) == 0 && engine.FollowRedirects { // If neither A nor AAAA records exist & redirects should be followed
		engine.DnsScanEngine, err = engine.DnsScanEngine.GetCNAMERecord(inputUrl) // Retrieve CNAME record if exists
		if err != nil {
			return engine, err
		}
		if engine.DnsScanEngine.CNAMERecord != "" { // If CNAME record exists
			if engine.Verbose {
				fmt.Println("No A or AAAA records for", inputUrl, ". Following CNAME...")
			}
			engine, err = engine.ScanDnsSimple(engine.dnsScanResult.CNAMERecord) // Follow CNAME recursively, but only checking A, AAAA, and CNAMEs
			if err != nil {
				return engine, err
			}
		}
	}

	return engine, nil
}
