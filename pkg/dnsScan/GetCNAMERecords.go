package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func GetCNAMERecord(url string, resolver *net.Resolver) (string, error) {
	var (
		err error

		record string
	)

	slog.Debug("dnsScan: Getting CNAME records started")

	record, err = resolver.LookupCNAME(context.Background(), url)
	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
		// No CNAME record available
	} else if err != nil {
		return record, err
	}

	slog.Debug("dnsScan: Getting CNAME records completed")

	return record, nil
}
