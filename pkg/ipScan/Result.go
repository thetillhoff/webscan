package ipScan

type Result struct {
	IpIsBlacklistedAt map[string][]string
	IpOwners          []string
	// BlacklistCheckUnavailable is true when at least one blacklist lookup could
	// not be completed (e.g. Spamhaus rejected the query). An empty
	// IpIsBlacklistedAt then does NOT mean the IPs are clean.
	BlacklistCheckUnavailable bool
}
