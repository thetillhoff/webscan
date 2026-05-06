package tlsScan

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func PrintResult(result Result, out io.Writer) {
	var (
		messages = []string{}
	)

	slog.Debug("tlsScan: Printing result started")

	if len(result.tlsScanResultPerIp) == 0 {
		slog.Debug("tlsScan: No results to print")
		return
	}

	// Shared certificate names
	sharedCertNames := result.ListSharedCertNames()
	if len(sharedCertNames) > 0 {
		messages = append(messages, "")
		messages = append(messages, "Certificate names (SN & SANs):")
		for _, certName := range sharedCertNames {
			messages = append(messages, fmt.Sprintf("- %s", certName))
		}
	}

	// Per-IP certificate name differences
	for ip := range result.tlsScanResultPerIp {
		nonSharedCertNames := result.ListNonSharedCertNamesForIp(ip)
		if len(nonSharedCertNames) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special certificate names on ip %s:", ip))
			for _, certName := range nonSharedCertNames {
				messages = append(messages, fmt.Sprintf("- %s", certName))
			}
		}
	}

	// Shared certificate issuers
	sharedCertIssuers := result.ListSharedCertIssuers()
	if len(sharedCertIssuers) > 0 {
		messages = append(messages, "")
		messages = append(messages, "Certificate issuers:")
		for _, certIssuer := range sharedCertIssuers {
			messages = append(messages, fmt.Sprintf("- %s", certIssuer))
		}
	}

	// Per-IP certificate issuer differences
	for ip := range result.tlsScanResultPerIp {
		nonSharedCertIssuers := result.ListNonSharedCertIssuersForIp(ip)
		if len(nonSharedCertIssuers) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special certificate issuers on ip %s:", ip))
			for _, certIssuer := range nonSharedCertIssuers {
				messages = append(messages, fmt.Sprintf("- %s", certIssuer))
			}
		}
	}

	// Shared TLS errors
	sharedTlsErrs := result.ListSharedTlsErr()
	if len(sharedTlsErrs) > 0 {
		messages = append(messages, "")
		messages = append(messages, "TLS errors:")
		for _, tlsErr := range sharedTlsErrs {
			if errors.Unwrap(tlsErr) != nil {
				messages = append(messages, errors.Unwrap(tlsErr).Error())
			} else {
				messages = append(messages, tlsErr.Error())
			}
			if strings.Contains(tlsErr.Error(), "connection reset by peer") {
				messages = append(messages, "This means the target might be available via TLS, but with a different hostname.")
			}
		}
	}

	// Per-IP TLS error differences
	for ip := range result.tlsScanResultPerIp {
		nonSharedTlsErrs := result.ListNonSharedTlsErrForIp(ip)
		if len(nonSharedTlsErrs) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special TLS errors on ip %s:", ip))
			for _, tlsErr := range nonSharedTlsErrs {
				if errors.Unwrap(tlsErr) != nil {
					messages = append(messages, errors.Unwrap(tlsErr).Error())
				} else {
					messages = append(messages, tlsErr.Error())
				}
				if strings.Contains(tlsErr.Error(), "connection reset by peer") {
					messages = append(messages, "This means the target might be available via TLS, but with a different hostname.")
				}
			}
		}
	}

	// Shared weak TLS versions
	sharedTlsVersions := result.ListSharedTlsVersions()
	weakTlsVersions := []uint16{}
	for _, tlsVersion := range sharedTlsVersions {
		if tlsVersion == tls.VersionTLS10 || tlsVersion == tls.VersionTLS11 {
			weakTlsVersions = append(weakTlsVersions, tlsVersion)
		}
	}
	if len(weakTlsVersions) > 0 {
		messages = append(messages, "")
		messages = append(messages, "Supported weak TLS versions:")
		for _, tlsVersion := range weakTlsVersions {
			switch tlsVersion {
			case tls.VersionTLS10:
				messages = append(messages, "- TLS 1.0")
			case tls.VersionTLS11:
				messages = append(messages, "- TLS 1.1")
			}
		}
	}

	// Per-IP weak TLS version differences
	for ip := range result.tlsScanResultPerIp {
		nonSharedTlsVersions := result.ListNonSharedTlsVersionsForIp(ip)
		weakVersions := []uint16{}
		for _, tlsVersion := range nonSharedTlsVersions {
			if tlsVersion == tls.VersionTLS10 || tlsVersion == tls.VersionTLS11 {
				weakVersions = append(weakVersions, tlsVersion)
			}
		}
		if len(weakVersions) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special supported weak TLS versions on ip %s:", ip))
			for _, tlsVersion := range weakVersions {
				switch tlsVersion {
				case tls.VersionTLS10:
					messages = append(messages, "- TLS 1.0")
				case tls.VersionTLS11:
					messages = append(messages, "- TLS 1.1")
				}
			}
		}
	}

	// Shared cipher rules
	sharedCipherRules := result.ListSharedCipherRules()
	if len(sharedCipherRules) > 0 {
		messages = append(messages, "")
		for rule, ciphers := range sharedCipherRules {
			messages = append(messages, "")
			messages = append(messages, rule)
			messages = append(messages, "Affected ciphers:")
			for _, affectedCipher := range ciphers {
				messages = append(messages, fmt.Sprintf("- %s", affectedCipher))
			}
		}
	}

	// Per-IP cipher rule differences
	for ip := range result.tlsScanResultPerIp {
		nonSharedCipherRules := result.ListNonSharedCipherRulesForIp(ip)
		if len(nonSharedCipherRules) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special cipher rules on ip %s:", ip))
			for rule, ciphers := range nonSharedCipherRules {
				messages = append(messages, "")
				messages = append(messages, rule)
				messages = append(messages, "Affected ciphers:")
				for _, affectedCipher := range ciphers {
					messages = append(messages, fmt.Sprintf("- %s", affectedCipher))
				}
			}
		}
	}

	if len(messages) > 0 {
		if _, err := fmt.Fprintf(out, "\n## TLS scan results\n"); err != nil {
			slog.Debug("tlsScan: Error writing to output", "error", err)
		}
		for _, message := range messages {
			if _, err := fmt.Fprintf(out, "%s\n", message); err != nil {
				slog.Debug("tlsScan: Error writing to output", "error", err)
			}
		}
	} else {
		slog.Debug("tlsScan: No information found")
	}

	slog.Debug("tlsScan: Printing result completed")
}
