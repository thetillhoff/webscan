package webscan

func (engine Engine) ScanDnsDetailed() (Engine, error) {
	var (
		err error
	)

	engine.DnsScanEngine, _ = engine.DnsScanEngine.GetDomainOwnerViaRDAP(engine.url)
	// // `err` is ignored here, as it's okay that it can't be retrieved. It's not a critical error, but an error nonetheless
	// if err != nil {
	// return engine, err
	// }

	engine.DnsScanEngine, err = engine.DnsScanEngine.GetAllRecords(engine.url)
	if err != nil {
		return engine, err
	}

	// TODO follow CNAME if exists

	return engine, nil
}
