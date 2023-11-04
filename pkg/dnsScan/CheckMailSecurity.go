package dnsScan

import "net"

func (engine Engine) CheckMailSecurity(url string, resolver *net.Resolver, dkimSelector string) []string {
	var (
		messages []string
		message  string
	)

	message = engine.CheckSpf()
	if message != "" {
		messages = append(messages, message)
	}

	message = engine.CheckDkim(dkimSelector+"._domainkey."+url, resolver)
	if message != "" {
		messages = append(messages, message)
	}

	message = engine.CheckDmarc(url, resolver)
	if message != "" {
		messages = append(messages, message)
	}

	return messages
}
