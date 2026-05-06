package ipScan

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	aRecords    []string
	aaaaRecords []string
	timeout     time.Duration
}

type ConfigOption = types.Option[scanConfig]

// WithARecords sets the A records to scan
func WithARecords(aRecords []string) ConfigOption {
	return func(sc *scanConfig) {
		sc.aRecords = aRecords
	}
}

// WithAAAARecords sets the AAAA records to scan
func WithAAAARecords(aaaaRecords []string) ConfigOption {
	return func(sc *scanConfig) {
		sc.aaaaRecords = aaaaRecords
	}
}

// WithTimeout sets the timeout for DNS and RDAP requests
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	var (
		result = Result{
			IpIsBlacklistedAt: map[string][]string{},
			IpOwners:          []string{},
		}

		maxIpAddressLength = 0
	)

	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	types.ApplyOptions(config, options)

	totalIPs := len(config.aRecords) + len(config.aaaaRecords)

	for _, aRecord := range config.aRecords {
		if len(aRecord) > maxIpAddressLength {
			maxIpAddressLength = len(aRecord)
		}
	}

	for _, aaaaRecord := range config.aaaaRecords {
		if len(aaaaRecord) > maxIpAddressLength {
			maxIpAddressLength = len(aaaaRecord)
		}
	}

	slog.Debug("ipScan: Scan started")

	switch {
	case totalIPs == 0:
		return result, errors.New("no ips to scan")
	case totalIPs == 1:
		slog.Debug("ipScan: Scanning one IP")
		status.SpinningUpdate("Scanning IP...") // Singular
	case totalIPs > 1:
		slog.Debug("ipScan: Scanning multiple IPs")
		status.SpinningXOfInit(totalIPs, "Scanning IPs...") // Plural
	}

	for _, aRecord := range config.aRecords {
		response, err := GetIPOwnerViaRDAP(aRecord, config.timeout)
		if err != nil {
			slog.Debug("ipScan: Error getting IP owner via RDAP for IPv4", "ipv4", aRecord, "error", err.Error())
			return result, err
		}
		result.IpOwners = append(result.IpOwners, fmt.Sprintf("According to RDAP information, IP %-*s is registered at %s", maxIpAddressLength, aRecord, response))

		blacklistMatches, err := IsIPBlacklisted(aRecord, config.timeout)
		if err != nil {
			slog.Debug("ipScan: Error on blacklist check of IPv4", "ipv4", aRecord, "error", err.Error())
			return result, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			result.IpIsBlacklistedAt[aRecord] = blacklistMatches
		}

		status.SpinningXOfUpdate()
	}

	for _, aaaaRecord := range config.aaaaRecords {
		response, err := GetIPOwnerViaRDAP(aaaaRecord, config.timeout)
		if err != nil {
			slog.Debug("ipScan: Error getting IP owner via RDAP for IPv6", "ipv6", aaaaRecord, "error", err.Error())
			return result, err
		}
		result.IpOwners = append(result.IpOwners, fmt.Sprintf("According to RDAP information, IP %-*s is registered at %s", maxIpAddressLength, aaaaRecord, response))

		blacklistMatches, err := IsIPBlacklisted(aaaaRecord, config.timeout)
		if err != nil {
			slog.Debug("ipScan: Error on blacklist check of IPv6", "ipv6", aaaaRecord, "error", err.Error())
			return result, err
		}

		if len(blacklistMatches) > 0 { // If ip was listed on at least one blacklist
			result.IpIsBlacklistedAt[aaaaRecord] = blacklistMatches
		}

		status.SpinningXOfUpdate()
	}

	switch {
	case totalIPs == 1:
		status.SpinningComplete("Scan of IP complete.") // Singular
	case totalIPs > 1:
		status.SpinningXOfComplete("Scan of IPs complete.") // Plural
	}

	slog.Debug("ipScan: Scan completed")

	return result, nil
}
