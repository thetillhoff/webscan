package dnsScan

func (engine Engine) CheckIpVersionCompatibility() Engine {
	var (
		message string = ""
	)

	if len(engine.ARecords) == 0 && len(engine.AAAARecords) == 0 {
		message = "No ips defined for domain"
	} else {
		if len(engine.AAAARecords) == 0 {
			message = "Hint: The resources of this domain are not reachable via IPv6."
		}

		if len(engine.ARecords) == 0 {
			message = "Hint: The resources of this domain are not reachable via IPv4."
		}
	}

	if message != "" {
		engine.OpinionatedHints = append(engine.OpinionatedHints, message)
	}

	return engine
}
