package dnsScan

import "net"

func (engine Engine) GetAllRecords(url string, resolver *net.Resolver) (Engine, error) {
	var (
		err error
	)

	// NS records
	engine, err = engine.GetNSRecords(url, resolver)
	if err != nil {
		return engine, err
	}

	// A records
	engine, err = engine.GetARecords(url, resolver)
	if err != nil {
		return engine, err
	}

	// AAAA records
	engine, err = engine.GetAAAARecords(url, resolver)
	if err != nil {
		return engine, err
	}

	// CNAME record
	engine, err = engine.GetCNAMERecord(url, resolver)
	if err != nil {
		return engine, err
	}

	// MX record
	engine, err = engine.GetMXRecords(url, resolver)
	if err != nil {
		return engine, err
	}

	// TXT record
	engine, err = engine.GetTXTRecords(url, resolver)
	if err != nil {
		return engine, err
	}

	return engine, nil
}
