package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetAAAARecords(url string) (Engine, error) {
	var (
		err error

		records []string
	)

	aaaaRecords, err := engine.resolver.LookupIP(context.Background(), "ip6", url)
	if _, ok := err.(*net.DNSError); ok {
		// No AAAA record available
		// TODO add opinionated recommendation to add it
	} else if err != nil {
		return engine, err
	}

	for _, ip := range aaaaRecords {
		records = append(records, ip.String())
	}

	engine.AAAARecords = records
	return engine, nil
}
