package dnsScan

import (
	"log/slog"
	"strings"

	"github.com/miekg/dns"
)

func CheckDkim(selectorUrl string, dnsClient *dns.Client, nameserver string) string {
	var (
		err         error
		txtRecords  []string
		cnameRecord string
	)

	slog.Debug("dnsScan: Checking dkim started")

	txtRecords, err = GetTXTRecords(selectorUrl, dnsClient, nameserver)
	if err != nil {
		// TODO add "Could not retrieve txt record" err/wrn/inf
		cnameRecord, err = GetCNAMERecord(selectorUrl, dnsClient, nameserver)
		if err != nil {
			return "Hint: Neither TXT nor CNAME found at specified DKIM selector."
		}

		return "Hint: DKIM selector redirects to " + cnameRecord
		// TODO recursively follow subDomainEngine.CNAMERecord
	}

	dkimRecord := ""
	for _, txtRecord := range txtRecords {
		if strings.HasPrefix(txtRecord, "v=DKIM1;") {
			if dkimRecord == "" { // Check if there was a dkim record detected before
				dkimRecord = txtRecord
			} else {
				return "Hint: Multiple DKIM records detected."
			}
		}
	}

	// TODO Verify dkimRecord

	slog.Debug("dnsScan: Checking dkim completed")

	return ""
}
