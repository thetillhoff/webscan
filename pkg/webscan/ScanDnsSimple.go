package webscan

import "fmt"

func (engine Engine) ScanDnsSimple() (Engine, error) {
	var (
		err error
	)

	fmt.Println("Scanning DNS (simple)...")

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetARecords(engine.url)
	if err != nil {
		return engine, err
	}

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetAAAARecords(engine.url)
	if err != nil {
		return engine, err
	}

	// TODO follow CNAME if exists

	return engine, nil
}
