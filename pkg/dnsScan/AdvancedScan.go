package dnsScan

import (
	"log/slog"

	"github.com/miekg/dns"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func AdvancedScan(status *status.Status, target types.Target, dnsClient *dns.Client, nameserver string, followRedirects bool) (Result, error) {
	var (
		err error

		result = Result{}
	)

	status.SpinningUpdate("Advanced scan of DNS running...")

	// `err` is ignored here, as it's okay that it can't be retrieved. It's not a critical error, but an error nonetheless
	result.DomainOwners, _ = GetDomainOwnerViaRDAP(target.Hostname())

	result.DomainIsBlacklistedAt, err = IsDomainBlacklisted(target.Hostname(), dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	result, err = GetAllRecords(target.Hostname(), dnsClient, nameserver)
	if err != nil {
		return result, err
	}

	if len(result.NSRecords) > 0 { // Only check nameserver owners if NS records exist
		result.NameserverOwners, err = GetNameserverOwnerViaRDAP(dnsClient, nameserver, result.NSRecords)
		if err != nil {
			return result, err
		}
	}

	if len(result.ARecords) == 0 && len(result.AAAARecords) == 0 && followRedirects { // If neither A nor AAAA records exist & redirects should be followed
		// TODO check if this should really go into record.CNAMERecord, or if for example a nested Result object would make more sense
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
				slog.Error("dnsScan: Failed to scan after following CNAME", "cname", result.CNAMERecord, "error", err)
				return result, err
			}
		}
	}

	// Domain Accessibility
	result.IpVersionCompatibility = CheckIpVersionCompatibility(result.ARecords, result.AAAARecords)
	result.DomainAccessibilityHints = GetDomainAccessibilityHints(target.Hostname())

	status.SpinningComplete("Advanced scan of DNS complete.")

	return result, nil
}
