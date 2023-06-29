package dnsScan

func (engine Engine) CheckMailSecurity(url string, dkimSelector string) []string {
	var (
		messages []string
		message  string
	)

	message = engine.CheckSpf()
	if message != "" {
		messages = append(messages, message)
	}

	message = engine.CheckDkim(dkimSelector + "._domainkey." + url)
	if message != "" {
		messages = append(messages, message)
	}

	message = engine.CheckDmarc(url)
	if message != "" {
		messages = append(messages, message)
	}

	return messages
}
