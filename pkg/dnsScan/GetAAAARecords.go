package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetAAAARecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for AAAA records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeAAAA)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No AAAA records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if aaaaRecord, ok := answer.(*dns.AAAA); ok {
			records = append(records, aaaaRecord.AAAA.String())
		}
	}

	slog.Debug("dnsScan: Checking for AAAA records completed", "url", url)

	return records, nil
}
