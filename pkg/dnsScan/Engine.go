package dnsScan

import (
	"context"
	"net"
	"time"
)

type Engine struct {
	// Internal variables
	resolver *net.Resolver

	// Scan Results
	NSRecords   []string
	ARecords    []string
	AAAARecords []string
	CNAMERecord string
	TXTRecords  []string
	MXRecords   []string
	SRVRecords  []string

	DomainOwners     []string
	OpinionatedHints []string
}

func DefaultEngine() Engine {
	return Engine{
		OpinionatedHints: []string{},

		resolver: nil, // Nil resolver is the same as a zero resolver

		NSRecords:   []string{},
		ARecords:    []string{},
		AAAARecords: []string{},
		CNAMERecord: "",
		TXTRecords:  []string{},

		DomainOwners: []string{},
	}
}

func EngineWithCustomDns(dnsServer string) Engine {
	engine := DefaultEngine()
	engine.resolver = &net.Resolver{
		PreferGo:     false,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
			return d.DialContext(ctx, network, net.JoinHostPort(dnsServer, "53"))
		},
	}
	return engine
}
