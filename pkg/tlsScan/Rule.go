package tlsScan

import "crypto/tls"

type Rule struct {
	description string
	matchFunc   func(cipherSuite tls.CipherSuite) bool
}
