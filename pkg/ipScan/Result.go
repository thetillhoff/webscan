package ipScan

type Result struct {
	IpIsBlacklistedAt map[string][]string
	IpOwners          []string
}
