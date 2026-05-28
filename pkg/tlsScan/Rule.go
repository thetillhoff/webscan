package tlsScan

import "crypto/tls"

type Rule struct {
	title       string
	description string
	matchFunc   func(cipherSuite tls.CipherSuite) bool
}
