package dnsScan

import (
	"net"
	"strings"
)

func (engine Engine) CheckDmarc(url string, resolver *net.Resolver) string {
	subDomainEngine := engine

	subDomainEngine, txtErr := subDomainEngine.GetTXTRecords("_dmarc."+url, resolver)
	if txtErr != nil {
		subDomainEngine, cnameErr := subDomainEngine.GetCNAMERecord("_dmarc."+url, resolver)
		if cnameErr != nil {
			return "Hint: Neither TXT nor CNAME records are set up for DMARC."
		}

		return "Hint: DKIM selector redirects to " + subDomainEngine.CNAMERecord
		// TODO recursively follow subDomainEngine.CNAMERecord
	}

	dmarcRecord := ""
	for _, txtRecord := range engine.TXTRecords {
		if strings.HasPrefix(txtRecord, "v=DMARC1;") {
			if dmarcRecord == "" { // Check if there was a dmarc record detected before
				dmarcRecord = txtRecord
			} else {
				return "Hint: Multiple DMARC records detected."
			}
		}
	}

	// TODO Verify dmarcRecord

	return ""
}
