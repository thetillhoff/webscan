package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetNSRecords(url string, resolver *net.Resolver) ([]string, error) {
	var (
		err error

		records   = []string{}
		nsRecords []*net.NS
	)

	slog.Debug("dnsScan: Getting MX records started")

	nsRecords, err = resolver.LookupNS(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No NS record available
	} else if err != nil {
		return records, err
	}

	for _, record := range nsRecords {
		records = append(records, record.Host)
	}

	slog.Debug("dnsScan: Getting MX records completed")

	return records, nil
}
