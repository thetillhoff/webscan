package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetSRVRecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for SRV records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeSRV)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No SRV records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if srvRecord, ok := answer.(*dns.SRV); ok {
			records = append(records, srvRecord.Target)
		}
	}

	slog.Debug("dnsScan: Checking for SRV records completed", "url", url)

	return records, nil
}
