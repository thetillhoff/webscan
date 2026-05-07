package dnsScan

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

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
	timeout         time.Duration
}

type ConfigOption = types.Option[scanConfig]

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

// WithTimeout sets the timeout for DNS queries
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

// resolveNameserver determines which nameserver to use based on configuration and system settings.
// It uses either an explicit custom nameserver or the system resolver configuration.
func resolveNameserver(customNameserver string) (string, *dns.ClientConfig, error) {
	defaultNameserver := "1.1.1.1:53"

	if customNameserver != "" {
		normalized := normalizeNameserver(customNameserver)
		slog.Debug("dnsScan: Using custom nameserver", "nameserver", normalized)
		return normalized, nil, nil
	}

	if runtime.GOOS == "windows" {
		slog.Debug("dnsScan: resolv.conf not available on Windows, using fallback", "fallback", defaultNameserver)
		return defaultNameserver, nil, nil
	}

	// Load system nameservers from resolv.conf
	resolvConfig, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return "", nil, fmt.Errorf("failed to load system DNS config: %w (set --dns to specify a resolver)", err)
	}

	if len(resolvConfig.Servers) > 0 {
		// Use the first nameserver as primary (like the system resolver does)
		server := resolvConfig.Servers[0]

		// Handle IPv6 addresses properly by wrapping them in square brackets
		var primaryNameserver string
		if types.IsIPv6(server) {
			// This is an IPv6 address, wrap it in square brackets
			primaryNameserver = "[" + server + "]:53"
		} else {
			// This is an IPv4 address or hostname, just add the port
			primaryNameserver = server + ":53"
		}

		slog.Debug("dnsScan: Using system nameservers from resolv.conf", "primary", primaryNameserver, "fallbacks", resolvConfig.Servers[1:])
		return primaryNameserver, resolvConfig, nil
	}

	return "", nil, errors.New("no DNS servers found in system config (set --dns to specify a resolver)")
}

func normalizeNameserver(input string) string {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return ""
	}

	host, port, err := net.SplitHostPort(raw)
	if err == nil {
		if port == "" {
			port = "53"
		}
		return net.JoinHostPort(host, port)
	}

	// If this is a plain IP/hostname without port, use default DNS port.
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		raw = strings.TrimPrefix(strings.TrimSuffix(raw, "]"), "[")
	}
	return net.JoinHostPort(raw, strconv.Itoa(53))
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

	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	types.ApplyOptions(config, options)

	slog.Debug("dnsScan: Scan started")
	status.SpinningUpdate("Scanning DNS records...")

	config.dnsClient = new(dns.Client)
	config.dnsClient.Net = "udp"
	config.dnsClient.Dialer = &net.Dialer{
		Timeout: config.timeout,
	}

	// Resolve which nameserver to use
	config.nameserver, config.resolvConfig, err = resolveNameserver(config.nameserver)
	if err != nil {
		slog.Error("dnsScan: could not determine nameserver", "error", err)
		status.SpinningComplete("Scan of DNS records failed.")
		return result, err
	}

	switch {
	case target.TargetType() == types.Domain && target.Schema() == types.NONE && config.advanced:
		slog.Info("dnsScan: Input identified as domain without schema (advanced)")

		result, err = AdvancedScan(
			status,
			target,
			config.dnsClient,
			config.nameserver,
			config.followRedirects,
			config.timeout,
		)
	case target.TargetType() == types.Domain:
		slog.Info("dnsScan: Input identified as domain", "schema", target.Schema().String())

		result, err = SimpleScan(
			target,
			config.dnsClient,
			config.nameserver,
			config.followRedirects,
		)
	case target.TargetType() == types.Ipv4:
		slog.Info("dnsScan: Input identified as IPv4")
		result.ARecords = []string{target.Hostname()}
	case target.TargetType() == types.Ipv6:
		slog.Info("dnsScan: Input identified as IPv6")
		result.AAAARecords = []string{target.Hostname()}
	default:
		slog.Error("dnsScan: Scan failed", "targetType", target.TargetType())
		status.SpinningComplete("Scan of DNS records failed.")
		return result, errors.ErrUnsupported // Unreachable code, since the type and the corresponding error handling happens earlier
	}

	if err != nil {
		slog.Error("dnsScan: Scan failed", "error", err)
		status.SpinningComplete("Scan of DNS records failed.")
		return result, err
	}

	slog.Debug("dnsScan: Scan completed")
	status.SpinningComplete("Scan of DNS records complete.")

	return result, nil
}
