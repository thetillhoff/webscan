package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetCNAMERecord(url string) (Engine, error) {
	var (
		err error

		record string
	)

	record, err = engine.resolver.LookupCNAME(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No CNAME record available
	} else if err != nil {
		return engine, err
	}

	engine.CNAMERecord = record
	return engine, nil
}
