package dnsScan

import (
	"errors"
	"log/slog"
	"net"

	"github.com/miekg/dns"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	nameserver      string
	dnsClient       *dns.Client
	resolvConfig    *dns.ClientConfig
	advanced        bool
	followRedirects bool
}

// ConfigOption represents a configuration option for DNS scanning
type ConfigOption func(*scanConfig)

// WithCustomNameServer sets a custom nameserver to use for DNS queries
func WithCustomNameServer(nameserver string) ConfigOption {
	return func(sc *scanConfig) {
		sc.nameserver = nameserver
	}
}

// WithAdvanced enables advanced scanning
func WithAdvanced(advanced bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.advanced = advanced
	}
}

// WithFollowRedirects enables following redirects
func WithFollowRedirects(followRedirects bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.followRedirects = followRedirects
	}
}

// isIPv6 checks if the given string is a valid IPv6 address
func isIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}

// resolveNameserver determines which nameserver to use based on configuration and system settings
func resolveNameserver(customNameserver string) (string, *dns.ClientConfig) {
	defaultNameserver := "1.1.1.1:53"

	if customNameserver != "" {
		slog.Debug("dnsScan: Using custom nameserver", "nameserver", customNameserver)
		return customNameserver, nil
	}

	// Load system nameservers from resolv.conf
	resolvConfig, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		slog.Debug("dnsScan: Failed to load resolv.conf, using fallback", "error", err, "fallback", defaultNameserver)
		return defaultNameserver, nil
	}

	if len(resolvConfig.Servers) > 0 {
		// Use the first nameserver as primary (like the system resolver does)
		server := resolvConfig.Servers[0]

		// Handle IPv6 addresses properly by wrapping them in square brackets
		var primaryNameserver string
		if isIPv6(server) {
			// This is an IPv6 address, wrap it in square brackets
			primaryNameserver = "[" + server + "]:53"
		} else {
			// This is an IPv4 address or hostname, just add the port
			primaryNameserver = server + ":53"
		}

		slog.Debug("dnsScan: Using system nameservers from resolv.conf", "primary", primaryNameserver, "fallbacks", resolvConfig.Servers[1:])
		return primaryNameserver, resolvConfig
	}

	slog.Debug("dnsScan: No nameservers found in resolv.conf, using fallback", "fallback", defaultNameserver)
	return defaultNameserver, nil
}

// Scans DNS records for the target.
//
// Custom resolver can be provided via WithResolver.
// Advanced scan can be enabled via WithAdvanced.
// Follow redirects can be enabled via WithFollowRedirects.
func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	var (
		err    error
		result Result
	)

	// Apply configuration options
	config := &scanConfig{}
	for _, option := range options {
		option(config)
	}

	slog.Debug("dnsScan: Scan started")

	// Initialize DNS client
	config.dnsClient = new(dns.Client)
	config.dnsClient.Net = "tcp" // Use TCP to handle large DNS responses
	config.dnsClient.Dialer = &net.Dialer{
		// Timeout:   200 * time.Millisecond, // TODO add timeout
	}

	// Resolve which nameserver to use
	config.nameserver, config.resolvConfig = resolveNameserver(config.nameserver)

	switch {
	case target.TargetType() == types.Domain && target.Schema() == types.NONE:
		slog.Info("input identified as domain without schema")

		result, err = AdvancedScan(
			status,
			target,
			config.dnsClient,
			config.nameserver,
			config.followRedirects,
		)
	case target.TargetType() == types.Domain && target.Schema() != types.NONE:
		slog.Info("input identified as domain with", "schema", target.Schema().String())

		result, err = SimpleScan(
			target,
			config.dnsClient,
			config.nameserver,
			config.followRedirects,
		)
	case target.TargetType() == types.Ipv4:
		slog.Info("input identified as ipv4")
		result.ARecords = []string{target.Hostname()}
	case target.TargetType() == types.Ipv6:
		slog.Info("input identified as ipv6")
		result.AAAARecords = []string{target.Hostname()}
	default:
		slog.Error("dnsScan: Scan failed", "targetType", target.TargetType())
		return result, errors.ErrUnsupported // Unreachable code, since the type and the corresponding error handling happens earlier
	}

	if err != nil {
		slog.Error("dnsScan: Scan failed", "error", err)
		return result, err
	}

	slog.Debug("dnsScan: Scan completed")

	return result, nil
}
