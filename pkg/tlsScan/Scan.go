package tlsScan

import (
	"errors"
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/portScan"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func Scan(target types.Target, status *status.Status, aRecords []string, aaaaRecords []string, portScanResult portScan.Result) (Result, error) {
	var (
		err    error
		result = Result{
			tlsScanResultPerIp: map[string]TlsScanResult{},
		}

		ips = append(aRecords, aaaaRecords...)

		ipsToScan = []string{}
	)

	slog.Debug("tlsScan: Scan started")

	if len(ips) == 0 {
		slog.Info("tlsScan: No ips to scan, skipping tls scan.")
		return result, nil
	}

	switch {
	case target.Schema() == types.HTTP:
		slog.Info("tlsScan: HTTP schema detected, skipping tls scan.")
		return result, nil
	case target.Port() == "":
		target.OverridePort("443")
		slog.Info("tlsScan: No schema detected, ignoring and using default port 443. If you want to set the schema to HTTP, and port to 80, add `:80` to the target.")
	case !portScanResult.IsPortOpen(target.PortAsUint16()):
		slog.Info("tlsScan: Port is closed, skipping tls scan.", "port", target.Port())
		return result, nil
	}

	for _, ip := range ips {
		if portScanResult.IsPortOpenOnIp(ip, target.PortAsUint16()) {
			ipsToScan = append(ipsToScan, ip)
		}
	}

	if len(ipsToScan) == 0 {
		switch len(ips) {
		case 1:
			slog.Info("tlsScan: Port is closed on target ip, skipping tls scan.", "port", target.Port())
		default:
			slog.Info("tlsScan: Port is closed on all target ips, skipping tls scan.", "port", target.Port())
		}
		return result, nil
	}

	for _, ip := range ipsToScan {
		tlsScanResult := TlsScanResult{}

		tlsScanResult.certInfos, tlsScanResult.tlsErr, err = evaluateTlsCertificate(target, ip) // Verify validity of tls certificate
		if err != nil {
			if errors.Unwrap(err) != nil {
				err = errors.Unwrap(err)
			}
			slog.Error("Could not evaluate tls certificate", "error", err)
			return result, nil // No need to continue if the tlsDial didn't work before, but returning an error ends the scan prematurely
		}
		tlsScanResult.enabledTlsVersions = scanEnabledTlsVersions(status, target, ip) // Try all tls versions
		tlsScanResult.enabledTlsCiphers = scanEnabledTlsCiphers(status, target, ip)   // Try all ciphers

		tlsScanResult.cipherRulesEvaluationResult = evaluateTlsCipherRules(tlsScanResult.enabledTlsCiphers) // Evaluate cipher rules

		result.tlsScanResultPerIp[ip] = tlsScanResult
	}

	slog.Debug("tlsScan: Scan completed")

	return result, nil
}
