package tlsScan

import (
	"crypto/tls"
	"fmt"
)

// Check available tls version and ciphers against best practices and print deviations
func PrintRecommendations(tlsCiphers []TlsCipher) {
	var (
		tls10allowed bool = false
		tls11allowed bool = false

		tlsWarnings    = []string{}
		cipherWarnings = []string{}
	)

	// Verify TLS version
	for _, tlsCipher := range tlsCiphers {
		if tlsCipher.TlsVersion == tls.VersionTLS10 { // Warn _once_ if tls 1.0 is allowed
			tls10allowed = true
			continue // Don't check cipher if TLS 1.0 is used for it
		} else if tlsCipher.TlsVersion == tls.VersionTLS11 { // Warn _once_ if tls 1.1 is allowed
			tls11allowed = true
			continue // Don't check cipher if TLS 1.1 is used for it
		}

		isCipherSecure := true
		for _, insecureCipher := range tls.InsecureCipherSuites() {
			if tlsCipher.Cipher == insecureCipher.ID {
				isCipherSecure = false
				break
			}
		}
		if !isCipherSecure {
			cipherWarnings = append(cipherWarnings, "Weak cipher allowed: "+tls.CipherSuiteName(tlsCipher.Cipher)) // TODO with Go1.21 it might be possible to print the tlsVersion along with it: https://github.com/golang/go/issues/46308
		}
	}

	if tls10allowed {
		tlsWarnings = append(tlsWarnings, "Weak tls version 1.0 allowed")
	}
	if tls11allowed {
		tlsWarnings = append(tlsWarnings, "Weak tls version 1.1 allowed")
	}

	if len(tlsWarnings) > 0 || len(cipherWarnings) > 0 {
		fmt.Println()
	}
	if len(tlsWarnings) > 0 {
		for _, warning := range tlsWarnings {
			fmt.Println(warning)
		}
	}
	if len(cipherWarnings) > 0 {
		for _, warning := range cipherWarnings {
			fmt.Println(warning)
		}
	}
	if len(tlsWarnings) > 0 || len(cipherWarnings) > 0 {
		fmt.Println("For recommendations for your webserver, visit https://ssl-config.mozilla.org/")
	}
}
