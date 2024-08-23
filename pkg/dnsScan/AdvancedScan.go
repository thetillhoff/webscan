package dnsScan

import (
	"log/slog"
	"net"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

// TODO put resolver always as first argument in all methods where it's needed

func AdvancedScan(status *status.Status, resolver *net.Resolver, inputUrl string, opinionated bool, followRedirects bool) (Result, error) {
	var (
		err error

		result = Result{}
	)

	slog.Debug("dnsScan: Advanced scan of " + inputUrl + " started")

	status.SpinningUpdate("Advanced scan of DNS running...")

	result.DomainOwners, _ = GetDomainOwnerViaRDAP(inputUrl)
	// // `err` is ignored here, as it's okay that it can't be retrieved. It's not a critical error, but an error nonetheless

	result.DomainIsBlacklistedAt, err = IsDomainBlacklisted(inputUrl, resolver)
	if err != nil {
		return result, err
	}

	result, err = GetAllRecords(inputUrl, resolver)
	if err != nil {
		return result, err
	}

	result.NameserverOwners, err = GetNameserverOwnerViaRDAP(resolver, result.NSRecords)
	if err != nil {
		return result, err
	}

	if opinionated {
		// Domain Accessibility
		result.IpVersionCompatibility = CheckIpVersionCompatibility(result.ARecords, result.AAAARecords) // TODO What if CNAME exists? See if function below
		result.DomainAccessibilityHints = GetDomainAccessibilityHints(inputUrl)
	}

	if len(result.ARecords) == 0 && len(result.AAAARecords) == 0 && followRedirects { // If neither A nor AAAA records exist & redirects should be followed
		// TODO check if this should really go into record.CNAMERecord, or if for example a nested Result object would make more sense
		result.CNAMERecord, err = GetCNAMERecord(inputUrl, resolver) // Retrieve CNAME record if exists
		if err != nil {
			return result, err
		}
		if result.CNAMERecord != "" { // If CNAME record exists
			// TODO DBG or INF
			slog.Warn("No A or AAAA records for", inputUrl, ". Following CNAME...")

			result, err = SimpleScan(resolver, result.CNAMERecord, followRedirects) // Follow CNAME recursively, but only checking A, AAAA, and CNAMEs
			if err != nil {
				return result, err
			}
		}
	}

	slog.Debug("dnsScan", "result", result)

	status.SpinningComplete("Advanced scan of DNS complete.")

	slog.Debug("dnsScan: Advanced scan completed")

	return result, nil
}
