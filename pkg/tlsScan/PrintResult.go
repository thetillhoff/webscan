package tlsScan

import (
	"crypto/tls"
	"fmt"
	"log/slog"
)

func PrintResult(result Result) {

	slog.Debug("tlsScan: Printing result started")

	//

	if result.tlsErr != nil || len(result.cipherRulesEvaluationResult) > 0 {

		fmt.Printf("\n\n## TLS scan results")

	}

	if len(result.CertNames) > 0 {
		fmt.Println()
		fmt.Println("Certificate names (SN first, then SANs):")
		for certName := range result.CertNames {
			fmt.Println("-", certName)
		}
	}

	if result.tlsErr != nil {
		fmt.Println()

		// Print tls certificate issues
		fmt.Println(result.tlsErr)
	}

	if len(result.cipherRulesEvaluationResult) > 0 {
		fmt.Println()

		// Print tls version issues
		for _, tlsVersion := range result.enabledTlsVersions {
			switch tlsVersion {
			case tls.VersionTLS10:
				fmt.Println("Weak tls version 1.0 is enabled.")
			case tls.VersionTLS11:
				fmt.Println("Weak tls version 1.1 is enabled.")
			}
		}
	}

	if len(result.cipherRulesEvaluationResult) > 0 {

		// // Print tls cipher issues
		for description, affectedCiphers := range result.cipherRulesEvaluationResult {
			fmt.Println()
			fmt.Println(description)
			fmt.Println("Affected ciphers:")
			for _, affectedCipher := range affectedCiphers {
				fmt.Println("-", affectedCipher)
			}
		}
	}

	slog.Debug("tlsScan: Printing result completed")
}
