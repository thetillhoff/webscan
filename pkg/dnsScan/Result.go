package dnsScan

type Result struct {
	NSRecords   []string
	ARecords    []string
	AAAARecords []string
	CNAMERecord string
	TXTRecords  []string
	MXRecords   []string
	SRVRecords  []string

	DomainOwners     []string
	NameserverOwners []string
	OpinionatedHints []string

	DomainIsBlacklistedAt []string

	IpVersionCompatibility   string
	DomainAccessibilityHints []string
}

// hasPrintableRecords reports whether PrintAllDnsRecords would emit any record.
func (result Result) hasPrintableRecords() bool {
	return len(result.NSRecords) > 0 ||
		len(result.ARecords) > 0 ||
		len(result.AAAARecords) > 0 ||
		result.CNAMERecord != "" ||
		len(result.MXRecords) > 0 ||
		len(result.TXTRecords) > 0
}
