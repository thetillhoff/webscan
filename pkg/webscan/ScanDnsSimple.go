package webscan

func (engine Engine) ScanDnsSimple() (Engine, error) {
	var (
		err error
	)

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
