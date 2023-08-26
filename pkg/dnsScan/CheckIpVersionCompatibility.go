package dnsScan

func (engine Engine) CheckIpVersionCompatibility(aRecords []string, aaaaRecords []string) Engine {
	var (
		message string = ""
	)

	if len(aRecords) == 0 && len(aaaaRecords) == 0 {
		message = "No ips defined for domain"
	} else {
		if len(aaaaRecords) == 0 {
			message = "Hint: The resources of this domain are not reachable via IPv6."
		}

		if len(aRecords) == 0 {
			message = "Hint: The resources of this domain are not reachable via IPv4."
		}
	}

	if message != "" {
		engine.OpinionatedHints = append(engine.OpinionatedHints, message)
	}

	return engine
}
