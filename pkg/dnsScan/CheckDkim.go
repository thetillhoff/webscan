package dnsScan

import (
	"net"
	"strings"
)

func (engine Engine) CheckDkim(selectorUrl string, resolver *net.Resolver) string {
	subDomainEngine := engine

	subDomainEngine, txtErr := subDomainEngine.GetTXTRecords(selectorUrl, resolver)
	if txtErr != nil {
		subDomainEngine, cnameErr := subDomainEngine.GetCNAMERecord(selectorUrl, resolver)
		if cnameErr != nil {
			return "Hint: Neither TXT nor CNAME found at specified DKIM selector."
		}

		return "Hint: DKIM selector redirects to " + subDomainEngine.CNAMERecord
		// TODO recursively follow subDomainEngine.CNAMERecord
	}

	dkimRecord := ""
	for _, txtRecord := range engine.TXTRecords {
		if strings.HasPrefix(txtRecord, "v=DKIM1;") {
			if dkimRecord == "" { // Check if there was a dkim record detected before
				dkimRecord = txtRecord
			} else {
				return "Hint: Multiple DKIM records detected."
			}
		}
	}

	// TODO Verify dkimRecord

	return ""
}
