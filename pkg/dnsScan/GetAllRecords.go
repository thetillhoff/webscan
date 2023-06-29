package dnsScan

func (engine Engine) GetAllRecords(url string) (Engine, error) {
	var (
		err error
	)

	// NS records
	engine, err = engine.GetNSRecords(url)
	if err != nil {
		return engine, err
	}

	// A records
	engine, err = engine.GetARecords(url)
	if err != nil {
		return engine, err
	}

	// AAAA records
	engine, err = engine.GetAAAARecords(url)
	if err != nil {
		return engine, err
	}

	// CNAME record
	engine, err = engine.GetCNAMERecord(url)
	if err != nil {
		return engine, err
	}

	// MX record
	engine, err = engine.GetMXRecords(url)
	if err != nil {
		return engine, err
	}

	// TXT record
	engine, err = engine.GetTXTRecords(url)
	if err != nil {
		return engine, err
	}

	return engine, nil
}
