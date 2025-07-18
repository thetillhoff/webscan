package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetTXTRecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for TXT records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeTXT)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No TXT records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if txtRecord, ok := answer.(*dns.TXT); ok {
			records = append(records, txtRecord.Txt...)
		}
	}

	slog.Debug("dnsScan: Checking for TXT records completed", "url", url)

	return records, nil
}
