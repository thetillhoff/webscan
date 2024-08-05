package dnsScan

import (
	"fmt"
	"log/slog"
)

func (result *Result) PrintAllDnsRecords() {
	var (
		record string
	)

	slog.Debug("dnsScan: Printing all dns records started")

	// NS records
	for _, record := range result.NSRecords {
		fmt.Println("NS", record)
	}

	// A records
	for _, record = range result.ARecords {
		fmt.Println("A", record)
	}

	// AAAA records
	for _, record = range result.AAAARecords {
		fmt.Println("AAAA", record)
	}

	// CNAME record
	if result.CNAMERecord != "" {
		fmt.Println("CNAME", result.CNAMERecord)
	}

	// MX record
	for _, record = range result.MXRecords {
		fmt.Println("MX", record)
	}

	// TXT record
	for _, record = range result.TXTRecords {
		fmt.Println("TXT", record)
	}

	slog.Debug("dnsScan: Printing all dns records completed")
}
