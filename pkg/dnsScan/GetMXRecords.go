package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetMXRecords(url string, resolver *net.Resolver) ([]string, error) {
	var (
		err error

		records   = []string{}
		mxRecords []*net.MX
	)

	slog.Debug("dnsScan: Getting MX records started")

	mxRecords, err = resolver.LookupMX(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No MX record available
	} else if err != nil {
		return records, err
	}

	for _, record := range mxRecords {
		records = append(records, record.Host)
	}

	slog.Debug("dnsScan: Getting MX records completed")

	return records, nil
}
