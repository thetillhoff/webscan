package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func SimpleScan(target types.Target, dnsClient *dns.Client, nameserver string, followRedirects bool) (Result, error) {
	var (
		err error

		result = Result{}
	)

	result.ARecords, err = GetARecords(target.Hostname(), dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result.AAAARecords, err = GetAAAARecords(target.Hostname(), dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	if len(result.ARecords) == 0 && len(result.AAAARecords) == 0 && followRedirects { // If neither A nor AAAA records exist & redirects should be followed
		result.CNAMERecord, err = GetCNAMERecord(target.Hostname(), dnsClient, nameserver) // Retrieve CNAME record if exists
		if err != nil {
			return result, err
		}
		if result.CNAMERecord != "" { // If CNAME record exists
			slog.Info("No A or AAAA records for", target.Hostname(), ". Following CNAME...")
			newTarget, err := types.NewTarget(result.CNAMERecord)
			if err != nil {
				slog.Error("dnsScan: Could not create new target from CNAME", "cname", result.CNAMERecord, "error", err)
				return result, err
			}
			result, err = SimpleScan(newTarget, dnsClient, nameserver, followRedirects) // Follow CNAME recursively, but only checking A, AAAA, and CNAMEs
			if err != nil {
				slog.Error("dnsScan: Could not follow CNAME", "cname", result.CNAMERecord, "error", err)
				return result, err
			}
		}
	}

	return result, nil
}
