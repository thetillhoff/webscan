package dnsScan

import (
	"fmt"
	"log/slog"
	"net"
)

func SimpleScan(resolver *net.Resolver, inputUrl string, followRedirects bool) (Result, error) {
	var (
		err error

		result = Result{}
	)

	slog.Debug("dnsScan: Simple scan started")

	fmt.Println("Scanning DNS (simple)...")

	result.ARecords, err = GetARecords(inputUrl, resolver)
	if err != nil {
		return result, err
	}

	result.AAAARecords, err = GetAAAARecords(inputUrl, resolver)
	if err != nil {
		return result, err
	}

	if len(result.ARecords) == 0 && len(result.AAAARecords) == 0 && followRedirects { // If neither A nor AAAA records exist & redirects should be followed
		result.CNAMERecord, err = GetCNAMERecord(inputUrl, resolver) // Retrieve CNAME record if exists
		if err != nil {
			return result, err
		}
		if result.CNAMERecord != "" { // If CNAME record exists
			// TODO INF or DBG
			fmt.Println("No A or AAAA records for", inputUrl, ". Following CNAME...")
			result, err = SimpleScan(resolver, result.CNAMERecord, followRedirects) // Follow CNAME recursively, but only checking A, AAAA, and CNAMEs
			if err != nil {
				return result, err
			}
		}
	}

	slog.Debug("dnsScan: Simple scan completed")

	return result, nil
}
