package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) GetMXRecords(url string, resolver *net.Resolver) (Engine, error) {
	var (
		err error

		records   []string
		mxRecords []*net.MX
	)

	mxRecords, err = resolver.LookupMX(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No MX record available
	} else if err != nil {
		return engine, err
	}

	for _, record := range mxRecords {
		records = append(records, record.Host)
	}

	engine.MXRecords = records
	return engine, nil
}
