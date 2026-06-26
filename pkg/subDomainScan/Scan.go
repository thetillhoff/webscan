package subDomainScan

import (
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

type scanConfig struct {
	certNames []string
	timeout   time.Duration
}

// ConfigOption represents a configuration option for subdomain scanning
type ConfigOption func(*scanConfig)

// WithCertNames sets the certificate names
func WithCertNames(certNames []string) ConfigOption {
	return func(sc *scanConfig) {
		sc.certNames = certNames
	}
}

// WithTimeout sets the HTTP request timeout
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) Result {
	var (
		result = Result{
			subdomainsFromTlsScan: map[string]struct{}{},
			subdomainsFromCrtSh:   map[string]struct{}{},
		}
	)

	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	for _, option := range options {
		option(config)
	}

	slog.Debug("subDomainScan: Scan started")

	status.SpinningUpdate("Scanning subdomains...")

	for _, subdomain := range config.certNames {

		switch {
		case subdomain == "": // If subdomain candidate is empty
			slog.Debug("subDomainScan: Skipping empty subdomain")
			continue
		case strings.Contains(subdomain, "="): // If subdomain candidate contains strings like `CN=X,O=Y,C=Z`
			slog.Debug("subDomainScan: Skipping subdomain containing '='", "subdomain", subdomain)
			continue
		case net.ParseIP(subdomain) != nil: // If subdomain candidate is an ip address
			slog.Debug("subDomainScan: Skipping ip address", "ip", subdomain)
			continue
		case subdomain == target.Hostname(): // If subdomain candidate is the same as the target domain
			slog.Debug("subDomainScan: Skipping subdomain that is the same as the target domain", "subdomain", subdomain)
			continue
		case !strings.HasSuffix(subdomain, "."+target.Hostname()): // If subdomain candidate is a completely different domain (e.g. notexample.com vs example.com)
			slog.Debug("subDomainScan: Skipping subdomain that is not a subdomain of the target domain", "subdomain", subdomain)
			continue
		}

		// TODO: Think about this later, since it adds insights on existing subdomains with the '*.', which is lost if it's removed
		// subdomain = strings.TrimPrefix(subdomain, "*.") // Remove wildcards, as they are invalid dns names, but might contain valid subdomains

		result.subdomainsFromTlsScan[subdomain] = struct{}{} // Add unique entries
	}

	var crtStatus crtShStatus
	result.subdomainsFromCrtSh, crtStatus = CheckCertLogs(target, config.timeout)
	result.crtShStatus = crtStatus

	status.SpinningComplete("Scan of subdomains complete.")

	slog.Debug("subDomainScan: Scan completed")

	return result
}
