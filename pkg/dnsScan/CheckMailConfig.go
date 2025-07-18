package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
)

func CheckMailConfig(url string, dnsClient *dns.Client, nameserver string, txtRecords []string, dkimSelector string) []string {
	var (
		messages []string
		message  string
	)

	slog.Debug("dnsScan: Checking mail started")

	message = CheckSpf(txtRecords)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDkim(dkimSelector+"._domainkey."+url, dnsClient, nameserver)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDmarc(url, dnsClient, nameserver)
	if message != "" {
		messages = append(messages, message)
	}

	slog.Debug("dnsScan: Checking mail completed")

	return messages
}
