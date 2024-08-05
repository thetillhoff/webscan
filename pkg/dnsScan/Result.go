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
