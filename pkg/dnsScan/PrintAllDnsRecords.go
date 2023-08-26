package dnsScan

import (
	"fmt"
)

func (engine Engine) PrintAllDnsRecords() {
	var (
		record string
	)

	// NS records
	for _, record := range engine.NSRecords {
		fmt.Println("NS", record)
	}

	// A records
	for _, record = range engine.ARecords {
		fmt.Println("A", record)
	}

	// AAAA records
	for _, record = range engine.AAAARecords {
		fmt.Println("AAAA", record)
	}

	// CNAME record
	if engine.CNAMERecord != "" {
		fmt.Println("CNAME", engine.CNAMERecord)
	}

	// MX record
	for _, record = range engine.MXRecords {
		fmt.Println("MX", record)
	}

	// TXT record
	for _, record = range engine.TXTRecords {
		fmt.Println("TXT", record)
	}
}
