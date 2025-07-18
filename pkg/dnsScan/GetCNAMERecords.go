package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetCNAMERecord(url string, dnsClient *dns.Client, nameserver string) (string, error) {
	var (
		record = ""
	)

	slog.Debug("dnsScan: Checking for CNAME record started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeCNAME)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No CNAME record found", "url", url, "error", err)
		return record, err
	}

	for _, answer := range response.Answer {
		if cnameRecord, ok := answer.(*dns.CNAME); ok {
			record = cnameRecord.Target
			break
		}
	}

	slog.Debug("dnsScan: Checking for CNAME record completed", "url", url)

	return record, nil
}
