package dnsScan

type Engine struct {
	// Scan Results
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
}

func DefaultEngine() Engine {
	return Engine{
		OpinionatedHints: []string{},

		NSRecords:   []string{},
		ARecords:    []string{},
		AAAARecords: []string{},
		CNAMERecord: "",
		TXTRecords:  []string{},

		DomainOwners:     []string{},
		NameserverOwners: []string{},
	}
}
