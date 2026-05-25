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

	message = CheckSPF(txtRecords)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDKIM(dkimSelector+"._domainkey."+url, dnsClient, nameserver)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDMARC(url, dnsClient, nameserver)
	if message != "" {
		messages = append(messages, message)
	}

	slog.Debug("dnsScan: Checking mail completed")

	return messages
}
