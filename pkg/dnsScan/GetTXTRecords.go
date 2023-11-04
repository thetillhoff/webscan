package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetTXTRecords(url string, resolver *net.Resolver) (Engine, error) {
	var (
		err error
	)

	records, err := resolver.LookupTXT(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No TXT record available
	} else if err != nil {
		return engine, err
	}

	engine.TXTRecords = records
	return engine, nil
}
