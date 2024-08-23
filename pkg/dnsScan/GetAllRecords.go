package dnsScan

import (
	"log/slog"
	"net"
)

func GetAllRecords(url string, resolver *net.Resolver) (Result, error) {
	var (
		err error

		result = Result{}
	)

	slog.Debug("dnsScan: Getting all records started")

	// NS records
	result.NSRecords, err = GetNSRecords(url, resolver)
	if err != nil {
		return result, err
	}

	// A records
	result.ARecords, err = GetARecords(url, resolver)
	if err != nil {
		return result, err
	}

	// AAAA records
	result.AAAARecords, err = GetAAAARecords(url, resolver)
	if err != nil {
		return result, err
	}

	// CNAME record
	result.CNAMERecord, err = GetCNAMERecord(url, resolver)
	if err != nil {
		return result, err
	}

	// MX record
	result.MXRecords, err = GetMXRecords(url, resolver)
	if err != nil {
		return result, err
	}

	// TXT record
	result.TXTRecords, err = GetTXTRecords(url, resolver)
	if err != nil {
		return result, err
	}

	slog.Debug("dnsScan: Getting all records completed")

	return result, nil
}
