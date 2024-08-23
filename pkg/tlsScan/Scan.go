package tlsScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

func Scan(status *status.Status, inputUrl string) (Result, error) {
	var (
		result = Result{}
	)

	slog.Debug("tlsScan: Scan started")

	result.CertNames, result.tlsErr = evaluateTlsCertificate(inputUrl)   // Verify validity of tls certificate
	result.enabledTlsVersions = scanEnabledTlsVersions(status, inputUrl) // Try all tls versions
	result.enabledTlsCiphers = scanEnabledTlsCiphers(status, inputUrl)   // Try all ciphers

	result.cipherRulesEvaluationResult = evaluateTlsCipherRules(result.enabledTlsCiphers) // Evaluate cipher rules

	slog.Debug("tlsScan: Scan completed")

	return result, nil
}
