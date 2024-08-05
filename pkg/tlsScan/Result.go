package tlsScan

import "crypto/tls"

type Result struct {
	CertNames          map[string]struct{}
	tlsErr             error
	enabledTlsVersions []uint16
	enabledTlsCiphers  []tls.CipherSuite

	cipherRulesEvaluationResult map[string][]string
}
