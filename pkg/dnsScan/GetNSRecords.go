package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetNSRecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for NS records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeNS)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No NS records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if nsRecord, ok := answer.(*dns.NS); ok {
			records = append(records, nsRecord.Ns)
		}
	}

	slog.Debug("dnsScan: Checking for NS records completed", "url", url)

	return records, nil
}
