package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetARecords(url string, resolver *net.Resolver) ([]string, error) {
	var (
		err error

		records = []string{}
	)

	slog.Debug("dnsScan: Getting A records started")

	aRecords, err := resolver.LookupIP(context.Background(), "ip4", url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No A record available
		slog.Debug("dnsScan: No ip4 address", "url", url)
	} else if err != nil {
		return records, err
	}

	for _, ip := range aRecords {
		records = append(records, ip.String())
	}

	slog.Debug("dnsScan: Getting A records completed")

	return records, nil
}
