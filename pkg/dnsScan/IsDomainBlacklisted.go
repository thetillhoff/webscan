package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func IsDomainBlacklisted(domain string, dnsClient *dns.Client, nameserver string) ([]string, error) {
	var (
		blacklists = []string{}
	)

	slog.Debug("dnsScan: Checking if domain is blacklisted started", "domain", domain)

	// List of DNSBLs to check
	dnsBlacklists := []string{
		"zen.spamhaus.org",
		"bl.spamcop.net",
		"dnsbl.sorbs.net",
		"b.barracudacentral.org",
		"bl.blocklist.de",
	}

	for _, dnsbl := range dnsBlacklists {
		query := domain + "." + dnsbl

		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(query), dns.TypeA)

		response, _, err := dnsClient.Exchange(m, nameserver)
		if err != nil {
			slog.Debug("dnsScan: DNSBL query failed", "dnsbl", dnsbl, "error", err)
			continue
		}

		if len(response.Answer) > 0 {
			blacklists = append(blacklists, dnsbl)
		}
	}

	slog.Debug("dnsScan: Checking if domain is blacklisted completed", "domain", domain, "blacklists", blacklists)

	return blacklists, nil
}
