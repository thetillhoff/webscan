package dnsScan

import (
	"log/slog"
	"strings"

	"github.com/miekg/dns"
)

func CheckDMARC(url string, dnsClient *dns.Client, nameserver string) string {
	var (
		err         error
		txtRecords  []string
		cnameRecord string
	)

	slog.Debug("dnsScan: Checking DMARC started")

	txtRecords, err = GetTXTRecords("_dmarc."+url, dnsClient, nameserver)
	if err != nil {
		cnameRecord, err = GetCNAMERecord("_dmarc."+url, dnsClient, nameserver)
		if err != nil {
			return "Hint: Neither TXT nor CNAME records are set up for DMARC."
		}

		return "Hint: DMARC record redirects to " + cnameRecord
		// TODO recursively follow subDomainEngine.CNAMERecord
	}

	dmarcRecord := ""
	for _, txtRecord := range txtRecords {
		trimmedRecord := strings.TrimSpace(txtRecord)
		if strings.HasPrefix(trimmedRecord, "v=DMARC1;") || strings.HasPrefix(trimmedRecord, "v=DMARC1 ;") {
			if dmarcRecord == "" { // Check if there was a dmarc record detected before
				dmarcRecord = txtRecord
			} else {
				return "Hint: Multiple DMARC records detected."
			}
		}
	}

	// TODO Verify dmarcRecord

	slog.Debug("dnsScan: Checking DMARC completed")

	return ""
}
