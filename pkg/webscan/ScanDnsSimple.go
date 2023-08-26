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

	// TODO follow CNAME if exists

	return engine, nil
}
