package tlsScan

import (
	"log/slog"
	"net/url"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

func Scan(status *status.Status, rawTarget string, parsedUrl *url.URL) (Result, error) {
	var (
		result = Result{}
	)

	slog.Debug("tlsScan: Scan started")

	result.CertNames, result.tlsErr = evaluateTlsCertificate(rawTarget, parsedUrl.Host) // Verify validity of tls certificate
	result.enabledTlsVersions = scanEnabledTlsVersions(status, rawTarget)               // Try all tls versions
	result.enabledTlsCiphers = scanEnabledTlsCiphers(status, rawTarget)                 // Try all ciphers

	result.cipherRulesEvaluationResult = evaluateTlsCipherRules(result.enabledTlsCiphers) // Evaluate cipher rules

	slog.Debug("tlsScan: Scan completed")

	return result, nil
}
