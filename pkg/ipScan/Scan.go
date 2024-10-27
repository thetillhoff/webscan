package ipScan

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/thetillhoff/webscan/pkg/status"
)

func Scan(status *status.Status, aRecords []string, aaaaRecords []string) (Result, error) {
	var (
		err error

		totalIPs = len(aRecords) + len(aaaaRecords)

		response         string
		blacklistMatches []string
		result           = Result{
			IpIsBlacklistedAt: map[string][]string{},
			IpOwners:          []string{},
		}
	)

	slog.Debug("ipScan: Scan started")

	// TODO logging
	if totalIPs == 0 { // No ips

		return result, errors.New("no ips to scan")

	} else if totalIPs == 1 {

		slog.Debug("Scanning one ip")
		status.SpinningUpdate("Scanning IP...") // Singular

	} else { // If there is more than one IP

		slog.Debug("Scanning more than one ip")
		status.SpinningXOfInit(len(aRecords)+len(aaaaRecords), "Scanning IPs...") // Plural

	}

	
	maxLength := 0
	for _, aRecord := range aRecords {
		if len(aRecord) > maxLength {
			maxLength = len(aRecord)
		}
	}
	for _, aaaaRecord := range aaaaRecords {
		if len(aaaaRecord) > maxLength {
			maxLength = len(aaaaRecord)
		}
	}

	for _, aRecord := range aRecords {
		response, err = GetIPOwnerViaRDAP(aRecord)
		if err != nil {
			slog.Debug("ipScan: error on getting ip owner via rdap for ip4 address", "ipv6", aRecord, "error", err.Error())
			return result, err
		}
		result.IpOwners = append(result.IpOwners, fmt.Sprintf("According to RDAP information, IP %-*s is registered at %s", maxLength, aRecord, response))

		blacklistMatches, err = IsIpBlacklisted(aRecord, true) // TODO check verbose bool
		if err != nil {
			slog.Debug("ipScan: error on blacklist check of ip4 address", "ipv4", aRecord, "error", err.Error())
			return result, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			result.IpIsBlacklistedAt[aRecord] = blacklistMatches
		}

		status.SpinningXOfUpdate()
	}

	for _, aaaaRecord := range aaaaRecords {
		response, err = GetIPOwnerViaRDAP(aaaaRecord)
		if err != nil {
			slog.Debug("ipScan: error on getting ip owner via rdap for ip6 address", "ipv6", aaaaRecord, "error", err.Error())
			return result, err
		}
		result.IpOwners = append(result.IpOwners, fmt.Sprintf("According to RDAP information, IP %-*s is registered at %s", maxLength, aaaaRecord, response))

		blacklistMatches, err = IsIpBlacklisted(aaaaRecord, true) // TODO check verbose bool
		if err != nil {
			slog.Debug("ipScan: error on blacklist check of ip6 address", "ipv6", aaaaRecord, "error", err.Error())
			return result, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			result.IpIsBlacklistedAt[aaaaRecord] = blacklistMatches
		}

		status.SpinningXOfUpdate()
	}

	if totalIPs > 1 { // If there is more than one IP
		status.SpinningXOfComplete("Scan of IPs complete.") // Plural
	} else {
		status.SpinningComplete("Scan of IP complete.") // Singular
	}

	slog.Debug("ipScan: Scan completed")

	return result, nil
}
