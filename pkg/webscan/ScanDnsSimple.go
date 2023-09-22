package webscan

import "fmt"

func (engine Engine) ScanDnsSimple(inputUrl string) (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning DNS (simple) of", inputUrl, "...")

	engine.dnsScanEngine, err = engine.dnsScanEngine.GetARecords(inputUrl)
	if err != nil {
		return engine, err
	}

	engine.dnsScanEngine, err = engine.dnsScanEngine.GetAAAARecords(inputUrl)
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
