package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetARecords(url string, resolver *net.Resolver) (Engine, error) {
	var (
		err error

		records []string
	)

	aRecords, err := resolver.LookupIP(context.Background(), "ip4", url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No A record available
	} else if err != nil {
		return engine, err
	}

	for _, ip := range aRecords {
		records = append(records, ip.String())
	}

	engine.ARecords = records
	return engine, nil
}
