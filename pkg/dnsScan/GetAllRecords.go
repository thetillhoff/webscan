package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetAllRecords(url string, dnsClient *dns.Client, nameserver string) (Result, error) {
	var (
		err error

		result = Result{}
	)

	slog.Debug("dnsScan: Getting all records started", "url", url)

	result.ARecords, err = GetARecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.AAAARecords, err = GetAAAARecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.CNAMERecord, err = GetCNAMERecord(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.MXRecords, err = GetMXRecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.NSRecords, err = GetNSRecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.TXTRecords, err = GetTXTRecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.SRVRecords, err = GetSRVRecords(url, dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	slog.Debug("dnsScan: Getting all records completed", "url", url)

	return result, nil
}
