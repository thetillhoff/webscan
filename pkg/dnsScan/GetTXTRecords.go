package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetTXTRecords(url string, resolver *net.Resolver) ([]string, error) {
	var (
		err error

		records []string
	)

	slog.Debug("dnsScan: Getting TXT records started")

	records, err = resolver.LookupTXT(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No TXT record available
	} else if err != nil {
		return records, err
	}

	slog.Debug("dnsScan: Getting TXT records completed")

	return records, nil
}
