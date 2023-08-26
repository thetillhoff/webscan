package webscan

import "fmt"

func (engine Engine) ScanDnsDetailed(inputUrl string) (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning DNS (detailed) of", inputUrl, "...")

	engine.DnsScanEngine, _ = engine.DnsScanEngine.GetDomainOwnerViaRDAP(inputUrl)
	// // `err` is ignored here, as it's okay that it can't be retrieved. It's not a critical error, but an error nonetheless
	// if err != nil {
	// return engine, err
	// }

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetAllRecords(inputUrl)
	if err != nil {
		return engine, err
	}

	// TODO follow CNAME if exists

	return engine, nil
}
