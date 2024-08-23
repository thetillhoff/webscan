package dnsScan

import (
	"context"
	"log/slog"
	"net"

	"github.com/thetillhoff/webscan/pkg/ipScan"
)

func GetNameserverOwnerViaRDAP(resolver *net.Resolver, nsRecords []string) ([]string, error) {
	var (
		err error

		aRecords []net.IP
		owner    string

		uniqueOwners = map[string]struct{}{}
		owners       = []string{}
	)

	slog.Debug("dnsScan: Getting nameserver owner via rdap started")

	for _, nameserver := range nsRecords {

		aRecords, err = resolver.LookupIP(context.Background(), "ip4", nameserver)
		if err != nil {
			return owners, err
		}

		for _, aRecord := range aRecords {
			owner, err = ipScan.GetIPOwnerViaRDAP(aRecord.String())
			if err != nil {
				return owners, err
			}
			uniqueOwners[owner] = struct{}{}
		}

	}

	for uniqueOwner := range uniqueOwners {
		owners = append(owners, uniqueOwner)
	}

	if len(owners) == 0 {
		slog.Info("Could not retrieve Nameserver Owner")
	}

	slog.Debug("dnsScan: Getting nameserver owner via rdap completed")

	return owners, nil
}
