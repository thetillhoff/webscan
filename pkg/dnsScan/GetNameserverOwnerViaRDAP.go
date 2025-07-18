package dnsScan

import (
	"fmt"
	"log/slog"
	"net"
	"slices"
	"strings"

	"github.com/miekg/dns"
	"github.com/thetillhoff/webscan/v3/pkg/ipScan"
)

func GetNameserverOwnerViaRDAP(dnsClient *dns.Client, nameserver string, nsRecords []string) ([]string, error) {
	var (
		owners = []string{}
		// Track which domains we've already tried to avoid duplicate RDAP lookups
		triedDomains = make(map[string]bool)
		// Track successful domain lookups to reuse results
		domainResults = make(map[string][]string)
	)

	slog.Debug("dnsScan: Getting nameserver owners via RDAP started", "nsRecords", nsRecords)

	for _, nsRecord := range nsRecords {
		// Remove trailing dot if present
		nsRecord = dns.Fqdn(nsRecord)

		// Check if the nameserver is an IP address or a hostname
		if net.ParseIP(nsRecord) != nil {
			// It's an IP address - use IP owner lookup
			owner, err := ipScan.GetIPOwnerViaRDAP(nsRecord)
			if err != nil {
				slog.Debug("dnsScan: Could not get IP nameserver owner via RDAP", "nsRecord", nsRecord, "error", err)
				continue
			}
			if owner != "" {
				owners = append(owners, owner)
			}
		} else {
			// It's a hostname - try progressively shorter domain parts with deduplication
			domainOwners, err := getDomainOwnerViaRDAPRecursive(nsRecord, triedDomains, domainResults)
			if err != nil {
				slog.Debug("dnsScan: Could not get domain nameserver owner via RDAP", "nsRecord", nsRecord, "error", err)
				continue
			}
			owners = append(owners, domainOwners...)
		}
	}

	slices.Sort(owners)
	owners = slices.Compact(owners)

	slog.Debug("dnsScan: Getting nameserver owners via RDAP completed", "owners", owners)

	return owners, nil
}

// getDomainOwnerViaRDAPRecursive tries to get domain owner by recursively trying shorter domain parts
// It avoids duplicate lookups by tracking tried domains and reusing successful results
func getDomainOwnerViaRDAPRecursive(hostname string, triedDomains map[string]bool, domainResults map[string][]string) ([]string, error) {
	// Remove trailing dot if present
	hostname = strings.TrimSuffix(hostname, ".")

	// Split by dots
	parts := strings.Split(hostname, ".")

	// Need at least 2 parts for a valid domain
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid hostname: %s", hostname)
	}

	// Try progressively shorter domain parts
	for i := 0; i < len(parts)-1; i++ {
		domain := strings.Join(parts[i:], ".")

		// Check if we already have a successful result for this domain
		if result, exists := domainResults[domain]; exists {
			slog.Debug("dnsScan: Reusing cached domain result", "hostname", hostname, "domain", domain, "owners", result)
			return result, nil
		}

		// Check if we've already tried this domain (and it failed)
		if triedDomains[domain] {
			slog.Debug("dnsScan: Skipping already tried domain", "hostname", hostname, "domain", domain)
			continue
		}

		// Mark this domain as tried
		triedDomains[domain] = true

		slog.Debug("dnsScan: Trying domain for RDAP lookup", "hostname", hostname, "domain", domain)

		// Perform the RDAP lookup
		owners, err := GetDomainOwnerViaRDAP(domain)
		if err == nil && len(owners) > 0 {
			// Cache the successful result
			domainResults[domain] = owners
			slog.Debug("dnsScan: Found domain owner via RDAP", "hostname", hostname, "domain", domain, "owners", owners)
			return owners, nil
		}

		// If this domain failed, continue to the next shorter one
		slog.Debug("dnsScan: Domain lookup failed, trying shorter domain", "domain", domain, "error", err)
	}

	return nil, fmt.Errorf("no valid domain found for hostname: %s", hostname)
}
