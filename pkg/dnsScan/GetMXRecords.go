package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func GetMXRecords(url string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		records = []string{}
	)

	slog.Debug("dnsScan: Checking for MX records started", "url", url)

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(url), dns.TypeMX)

	response, _, err := dnsClient.Exchange(m, nameserver)
	if err != nil {
		slog.Debug("dnsScan: No MX records found", "url", url, "error", err)
		return records, err
	}

	for _, answer := range response.Answer {
		if mxRecord, ok := answer.(*dns.MX); ok {
			records = append(records, mxRecord.Mx)
		}
	}

	slog.Debug("dnsScan: Checking for MX records completed", "url", url)

	return records, nil
}
