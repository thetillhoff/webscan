package ipScan

func CheckIpVersionCompatibility(aRecords []string, aaaaRecords []string) string {

	if len(aRecords) == 0 && len(aaaaRecords) == 0 {
		return "No ips defined for domain"
	} else {
		if len(aaaaRecords) == 0 {
			return "Hint: The resources of this domain are not reachable via IPv6."
		}

		if len(aRecords) == 0 {
			return "Hint: The resources of this domain are not reachable via IPv4."
		}
	}

	return "" // No issues found
}
