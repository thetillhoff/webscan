package tlsScan

import (
	"crypto/tls"
	"strings"
)

// Check available tls version and ciphers against best practices and print deviations
func GetRecommendations(tlsCiphers []TlsCipher) []string {
	var (
		tls10allowed bool = false
		tls11allowed bool = false

		tlsWarnings    = []string{}
		cipherWarnings []string
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
	}

	if tls10allowed {
		tlsWarnings = append(tlsWarnings, "Weak tls version 1.0 allowed")
	}
	if tls11allowed {
		tlsWarnings = append(tlsWarnings, "Weak tls version 1.1 allowed")
	}

	// Verify ciphers (https://ciphersuite.info/cs/?tls=tls12&singlepage=true has some nice hints on the reasons behind deeming a cipher insecure)
	for _, tlsCipher := range tlsCiphers {

		cipherWarnings = []string{}

		isCipherSecure := true
		for _, insecureCipher := range tls.InsecureCipherSuites() {
			if tlsCipher.Cipher == insecureCipher.ID {
				isCipherSecure = false
				break
			}
		}

		// Note as of 2024-01:
		// While true, this is generally not considered to be a dealbreaker in TLS encryption.
		// if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "RSA") {
		// 	// RSA makes it possible to decrypt if the certificate is compromised in the future -> use ECDHE instead
		// 	cipherWarnings = append(cipherWarnings, "Recommending against RSA, as it's possible to decrypt traffic at a later time should the certificate be compromised in the future.")
		// }

		if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "3DES") {
			// 3DES has 64-bit blocks, which makes it fundamentally vulnerable to birthday attacks given enough traffic https://sweet32.info/
			cipherWarnings = append(cipherWarnings, "Recommending against 3DES, as it's vulnerable to birthday attacks given enough traffic (https://sweet32.info/).")
		}

		if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "RC4") {
			// RC4 has practically exploitable biases that can lead to plaintext recovery without side channels https://www.rc4nomore.com/ & https://blog.cloudflare.com/killing-rc4-the-long-goodbye/ & https://blog.cloudflare.com/end-of-the-road-for-rc4/ & https://datatracker.ietf.org/doc/html/rfc7465
			cipherWarnings = append(cipherWarnings, "Recommending against RC4, as it was prohibited to use by the IETF in 2015 (https://www.rc4nomore.com/, https://datatracker.ietf.org/doc/html/rfc7465).")
		}

		// Note as of 2024-01:
		// Copied from https://security.stackexchange.com/a/207414:
		// - Amazon s2n thought they fixed it, but turns out they didn't: https://eprint.iacr.org/2015/1129
		// - OpenSSL introduced a much worse vulnerability when they tried to fix it: https://www.openssl.org/news/secadv/20160503.txt
		// - Google's Adam Langley, possibly the best TLS implementer in the world, chose to not implement the fix in the Go standard library's TLS implementation and recommended people don't support CBC cipher suites if they're worried: https://twitter.com/agl__/status/669182140244824064
		// - The correct implementation of TLS CBC ciphersuites is much too difficult: https://www.imperialviolet.org/2013/02/04/luckythirteen.html
		// - Implementations thought fully patched and secure are discovered to be insecure as variations on the attack improve: https://eprint.iacr.org/2018/747
		// -> There is no _known_ attack for CBC, but it is considered weak and thereby insecure anyway.
		if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "CBC") {
			// CBC only with Encrypt-then-MAC -> recommend against it, as it's hard to get right https://blog.cloudflare.com/yet-another-padding-oracle-in-openssl-cbc-ciphersuites/
			cipherWarnings = append(cipherWarnings, "Recommending against CBC, as it seems fundamentally flawed since the Lucky13 vulnerability was discovered (https://en.wikipedia.org/wiki/Lucky_Thirteen_attack, https://security.stackexchange.com/a/207414).")
		}

		if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "ECDH_") {
			// ECDH_ only until 2026
			cipherWarnings = append(cipherWarnings, "Keep in mind ECDH_ ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.")
		}
		if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "DH_") {
			// DH_ only until 2026
			cipherWarnings = append(cipherWarnings, "Keep in mind DH_ ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.")
		}

		// Note as of 2024-01:
		// AES-128 is preferred over AES-256 as the larger key doesn't have any serious advantages in web traffic
		// But there are also no serious drawbacks, so a recommendation feels unnecessary

		// SHA256 or higher with CBC and GCM - CCM uses SHA256 anyway
		// TODO

		if !isCipherSecure || len(cipherWarnings) > 0 { // Prepend cipher name for warnings
			cipherWarnings = append([]string{"\nWeak cipher allowed: " + tls.CipherSuiteName(tlsCipher.Cipher)}, cipherWarnings...) // TODO with Go1.21 it might be possible to print the tlsVersion along with it: https://github.com/golang/go/issues/46308
		}

		tlsWarnings = append(tlsWarnings, cipherWarnings...)
	}

	if len(tlsWarnings) > 0 {
		tlsWarnings = append(tlsWarnings, "\nFor recommendations for your webserver, visit https://ssl-config.mozilla.org/")
	}

	return tlsWarnings
}
