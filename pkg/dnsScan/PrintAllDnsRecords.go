package dnsScan

import (
	"fmt"
	"io"
	"log/slog"
)

type record struct {
	recordTypeName string
	recordValue    string
}

func (result *Result) PrintAllDnsRecords(out io.Writer) {
	var (
		records = []record{}

		maxRecordTypeNameLength = 0
	)

	slog.Debug("dnsScan: Printing all dns records started")

	// NS records
	for _, recordValue := range result.NSRecords {
		recordTypeName := "NS"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: recordValue})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	// A records
	for _, recordValue := range result.ARecords {
		recordTypeName := "A"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: recordValue})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	// AAAA records
	for _, recordValue := range result.AAAARecords {
		recordTypeName := "AAAA"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: recordValue})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	// CNAME record
	if result.CNAMERecord != "" {
		recordTypeName := "CNAME"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: result.CNAMERecord})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	// MX record
	for _, recordValue := range result.MXRecords {
		recordTypeName := "MX"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: recordValue})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	// TXT record
	for _, recordValue := range result.TXTRecords {
		recordTypeName := "TXT"
		records = append(records, record{recordTypeName: recordTypeName, recordValue: recordValue})
		if len(recordTypeName) > maxRecordTypeNameLength {
			maxRecordTypeNameLength = len(recordTypeName)
		}
	}

	for _, record := range records {
		if _, err := fmt.Fprintf(out, "%-*s %s\n", maxRecordTypeNameLength, record.recordTypeName, record.recordValue); err != nil {
			slog.Debug("dnsScan: Error writing to output", "error", err)
		}
	}

	slog.Debug("dnsScan: Printing all dns records completed")
}
