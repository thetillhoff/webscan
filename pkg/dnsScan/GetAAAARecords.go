package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetAAAARecords(url string, resolver *net.Resolver) ([]string, error) {
	var (
		err error

		records = []string{}
	)

	slog.Debug("dnsScan: Getting AAAA records started")

	aaaaRecords, err := resolver.LookupIP(context.Background(), "ip6", url)
	if _, ok := err.(*net.DNSError); ok {
		// No AAAA record available
		slog.Debug("dnsScan: No ip6 address", "url", url)
	} else if err != nil {
		return records, err
	}

	for _, ip := range aaaaRecords {
		records = append(records, ip.String())
	}

	slog.Debug("dnsScan: Getting AAAA records completed")

	return records, nil
}
