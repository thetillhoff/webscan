package dnsScan

import "log/slog"

func CheckIpVersionCompatibility(aRecords []string, aaaaRecords []string) string {
	var (
		message string = ""
	)

	slog.Debug("dnsScan: Checking ip version compatibility started")

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

	slog.Debug("dnsScan: Checking ip version compatibility completed")

	return message
}
