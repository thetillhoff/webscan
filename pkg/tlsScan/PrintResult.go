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

	// Print tls certificate issues

	sharedCertNames := result.ListSharedCertNames()
	if len(sharedCertNames) > 0 {
		messages = append(messages, "")
		messages = append(messages, "Certificate names (SN & SANs):")
		for _, certName := range sharedCertNames {
			messages = append(messages, fmt.Sprintf("- %s", certName))
		}
	}

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

	sharedCertIssuers := result.ListSharedCertIssuers()
	if len(sharedCertIssuers) > 0 {
		messages = append(messages, "")
		messages = append(messages, "Certificate issuers:")
		for _, certIssuer := range sharedCertIssuers {
			messages = append(messages, fmt.Sprintf("- %s", certIssuer))
		}
	}

	for ip := range result.tlsScanResultPerIp {
		nonSharedCertIssuers := result.ListNonSharedCertIssuersForIp(ip)

		if len(nonSharedCertIssuers) > 0 {
			messages = append(messages, "")
			messages = append(messages, fmt.Sprintf("Special certificate issuers on ip %s:", ip))
			for _, certIssuer := range nonSharedCertIssuers {
				messages = append(messages, fmt.Sprintf("- %s", certIssuer))
			}
		}

		sharedTlsErrs := result.ListSharedTlsErr()
		if len(sharedTlsErrs) > 0 {
			messages = append(messages, "")
			messages = append(messages, "Tls errors:")
			for _, tlsErr := range sharedTlsErrs {
				if errors.Unwrap(tlsErr) != nil {
					messages = append(messages, errors.Unwrap(tlsErr).Error())
				} else {
					messages = append(messages, tlsErr.Error())
				}

				if strings.Contains(tlsErr.Error(), "connection reset by peer") {
					messages = append(messages, "This means the target might be available via TLS, but with a different hostname.")
				} else {
					slog.Debug("tlsScan: tlsErr is not a connection reset by peer error", "tlsErr", tlsErr)
				}
			}
		}

		for ip := range result.tlsScanResultPerIp {
			nonSharedTlsErrs := result.ListNonSharedTlsErrForIp(ip)
			if len(nonSharedTlsErrs) > 0 {
				messages = append(messages, "")
				messages = append(messages, fmt.Sprintf("Special tls errors on ip %s:", ip))

				for _, tlsErr := range nonSharedTlsErrs {
					if errors.Unwrap(tlsErr) != nil {
						messages = append(messages, errors.Unwrap(tlsErr).Error())
					} else {
						messages = append(messages, tlsErr.Error())
					}

					if strings.Contains(tlsErr.Error(), "connection reset by peer") {
						messages = append(messages, "This means the target might be available via TLS, but with a different hostname.")
					} else {
						slog.Debug("tlsScan: tlsErr is not a connection reset by peer error", "tlsErr", tlsErr)
					}
				}
			}
		}

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

		for ip := range result.tlsScanResultPerIp {
			nonSharedTlsVersions := result.ListNonSharedTlsVersionsForIp(ip)
			weakTlsVersions := []uint16{}
			for _, tlsVersion := range nonSharedTlsVersions {
				if tlsVersion == tls.VersionTLS10 || tlsVersion == tls.VersionTLS11 {
					weakTlsVersions = append(weakTlsVersions, tlsVersion)
				}
			}
			if len(weakTlsVersions) > 0 {
				messages = append(messages, "")
				messages = append(messages, fmt.Sprintf("Special supported weak TLS versions on ip %s:", ip))
				for _, tlsVersion := range weakTlsVersions {
					switch tlsVersion {
					case tls.VersionTLS10:
						messages = append(messages, "- TLS 1.0")
					case tls.VersionTLS11:
						messages = append(messages, "- TLS 1.1")
					}
				}
			}
		}

		// TODO: Improve the way the version issues are printed above, by taking into consideration the below.
		// // Print tls version issues
		// if len(result.enabledTlsVersions) > 0 {

		// 	weakTlsVersions := []string{}

		// 	for _, tlsVersion := range result.enabledTlsVersions {
		// 		switch tlsVersion {
		// 		case tls.VersionTLS10:
		// 			weakTlsVersions = append(weakTlsVersions, "TLS 1.0")
		// 		case tls.VersionTLS11:
		// 			weakTlsVersions = append(weakTlsVersions, "TLS 1.1")
		// 		}
		// 	}

		// 	if len(weakTlsVersions) > 0 {
		// 		messages = append(messages, "Weak tls versions are enabled:")
		// 		for _, weakTlsVersion := range weakTlsVersions {
		// 			messages = append(messages, fmt.Sprintf("- %s", weakTlsVersion))
		// 		}
		// 	}
		// }

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
}
