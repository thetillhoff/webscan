package dnsScan

import (
	"context"
	"net"

	"github.com/thetillhoff/webscan/pkg/ipScan"
)

func (engine Engine) GetNameserverOwnerViaRDAP(resolver *net.Resolver) (Engine, error) {
	var (
		err error

		aRecords []net.IP
		owner    string

		uniqueOwners = map[string]struct{}{}
		owners       = []string{}
	)

	for _, nameserver := range engine.NSRecords {

		aRecords, err = resolver.LookupIP(context.Background(), "ip4", nameserver)
		if err != nil {
			return engine, err
		}

		for _, aRecord := range aRecords {
			owner, err = ipScan.GetIPOwnerViaRDAP(aRecord.String())
			if err != nil {
				return engine, err
			}
			uniqueOwners[owner] = struct{}{}
		}

	}

	for uniqueOwner := range uniqueOwners {
		owners = append(owners, uniqueOwner)
	}
	engine.NameserverOwners = owners

	return engine, nil
}
