package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetNSRecords(url string) (Engine, error) {
	var (
		err error

		records   []string
		nsRecords []*net.NS
	)

	nsRecords, err = engine.resolver.LookupNS(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No NS record available
	} else if err != nil {
		return engine, err
	}

	for _, record := range nsRecords {
		records = append(records, record.Host)
	}

	engine.NSRecords = records
	return engine, nil
}
