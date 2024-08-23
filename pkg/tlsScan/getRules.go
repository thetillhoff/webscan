package tlsScan

import (
	"crypto/tls"
	"strings"
)

func getRules() []Rule {
	var (
		rules = []Rule{}
	)

	rules = append(rules, Rule{

		// Verify ciphers (https://ciphersuite.info/cs/?tls=tls12&singlepage=true has some nice hints on the reasons behind deeming a cipher insecure)
		description: `Some ciphers are deemed unsecure by Golang.
More information: https://ciphersuite.info/cs/?tls=tls12&singlepage=true`,
		matchFunc: func(cipherSuite tls.CipherSuite) bool {
			isCipherSecure := true
			for _, insecureCipher := range tls.InsecureCipherSuites() {
				if cipherSuite.ID == insecureCipher.ID {
					isCipherSecure = false
					break
				}
			}
			return isCipherSecure
		},
	})

	// Note as of 2024-01:
	// While true, this is generally not considered to be a dealbreaker in TLS encryption.
	// if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "RSA") {
	// 	// RSA makes it possible to decrypt if the certificate is compromised in the future -> use ECDHE instead
	// 	cipherWarnings = append(cipherWarnings, "Recommending against RSA, as it's possible to decrypt traffic at a later time should the certificate be compromised in the future.")
	// }

	// 3DES has 64-bit blocks, which makes it fundamentally vulnerable to birthday attacks given enough traffic https://sweet32.info/
	rules = append(rules, Rule{
		description: `3DES is fundamentally vulnerable to birthday attacks given enough traffic.
More information: https://sweet32.info/`,
		matchFunc: func(cipherSuite tls.CipherSuite) bool {
			return strings.Contains(cipherSuite.Name, "3DES")
		},
	})

	// RC4 has practically exploitable biases that can lead to plaintext recovery without side channels
	// https://www.rc4nomore.com/ & https://blog.cloudflare.com/killing-rc4-the-long-goodbye/ & https://blog.cloudflare.com/end-of-the-road-for-rc4/ & https://datatracker.ietf.org/doc/html/rfc7465
	rules = append(rules, Rule{
		description: `RC4 was prohibited by the IETF in 2015.
More information: https://www.rc4nomore.com/, https://datatracker.ietf.org/doc/html/rfc7465`,
		matchFunc: func(cipherSuite tls.CipherSuite) bool {
			return strings.Contains(cipherSuite.Name, "RC4")
		},
	})

	// Note as of 2024-01:
	// Copied from https://security.stackexchange.com/a/207414:
	// - Amazon s2n thought they fixed it, but turns out they didn't: https://eprint.iacr.org/2015/1129
	// - OpenSSL introduced a much worse vulnerability when they tried to fix it: https://www.openssl.org/news/secadv/20160503.txt
	// - Google's Adam Langley, possibly the best TLS implementer in the world, chose to not implement the fix in the Go standard library's TLS implementation and recommended people don't support CBC cipher suites if they're worried: https://twitter.com/agl__/status/669182140244824064
	// - The correct implementation of TLS CBC ciphersuites is much too difficult: https://www.imperialviolet.org/2013/02/04/luckythirteen.html
	// - Implementations thought fully patched and secure are discovered to be insecure as variations on the attack improve: https://eprint.iacr.org/2018/747
	// -> There is no _known_ attack for CBC, but it is considered weak and thereby insecure anyway.
	//
	// CBC only with Encrypt-then-MAC -> recommend against it, as it's hard to get right https://blog.cloudflare.com/yet-another-padding-oracle-in-openssl-cbc-ciphersuites/
	rules = append(rules, Rule{
		description: `CBC seems to be fundamentally flawed since the Lucky13 vulnerability was discovered - even though no known attacks exists. Yet.
More information: https://en.wikipedia.org/wiki/Lucky_Thirteen_attack, https://security.stackexchange.com/a/207414`,
		matchFunc: func(cipherSuite tls.CipherSuite) bool {
			return strings.Contains(cipherSuite.Name, "CBC")
		},
	})

	// TODO start displaying after 2026
	// if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "ECDH_") {
	// 	// ECDH_ only until 2026
	// 	cipherWarnings = append(cipherWarnings, "Keep in mind ECDH_ ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.")
	// }
	// if strings.Contains(tls.CipherSuiteName(tlsCipher.Cipher), "DH_") {
	// 	// DH_ only until 2026
	// 	cipherWarnings = append(cipherWarnings, "Keep in mind DH_ ciphers don't support Perfect Forward Secrecy and shouldn't be used after 2026.")
	// }

	// Note as of 2024-01:
	// AES-128 is preferred over AES-256 as the larger key doesn't have any serious advantages in web traffic
	// But there are also no serious drawbacks, so a recommendation feels unnecessary

	// SHA256 or higher with CBC and GCM - CCM uses SHA256 anyway

	return rules
}
