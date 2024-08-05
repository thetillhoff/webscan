package dnsScan

import (
	"log/slog"
	"net"
)

func CheckMailConfig(url string, resolver *net.Resolver, txtRecords []string, dkimSelector string) []string {
	var (
		messages []string
		message  string
	)

	slog.Debug("dnsScan: Checking mail started")

	message = CheckSpf(txtRecords)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDkim(dkimSelector+"._domainkey."+url, resolver)
	if message != "" {
		messages = append(messages, message)
	}

	message = CheckDmarc(url, resolver)
	if message != "" {
		messages = append(messages, message)
	}

	slog.Debug("dnsScan: Checking mail completed")

	return messages
}
