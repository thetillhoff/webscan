package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetARecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for A records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeA)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No A records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if aRecord, ok := answer.(*dns.A); ok {
			records = append(records, aRecord.A.String())
		}
	}

	slog.Debug("dnsScan: Checking for A records completed", "url", url)

	return records, nil
}
