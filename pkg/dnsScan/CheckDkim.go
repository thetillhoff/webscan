package dnsScan

import (
	"log/slog"
	"net"
	"strings"
)

func CheckDkim(selectorUrl string, resolver *net.Resolver) string {
	var (
		err         error
		txtRecords  []string
		cnameRecord string
	)

	slog.Debug("dnsScan: Checking dkim started")

	txtRecords, err = GetTXTRecords(selectorUrl, resolver)
	if err != nil {
		// TODO add "Could not retrieve txt record" err/wrn/inf
		cnameRecord, err = GetCNAMERecord(selectorUrl, resolver)
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
