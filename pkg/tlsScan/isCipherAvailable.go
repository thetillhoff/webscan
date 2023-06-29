package tlsScan

import (
	"crypto/tls"
	"os"
)

func isCipherAvailable(url string, tlsCipher TlsCipher, allowedCiphers chan<- TlsCipher) {
	defer wg.Done()

	_, err := tls.Dial("tcp", url+":443", &tls.Config{
		MinVersion:               tlsCipher.TlsVersion,
		MaxVersion:               tlsCipher.TlsVersion,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		ServerName:               url,
		CipherSuites: []uint16{
			tlsCipher.Cipher,
		},
	})

	// TODO try each TLS version
	// TODO try each cipher, warn on insecure, at least one secure

	if !os.IsTimeout(err) && err == nil { // If no timeout error occurred and there was no other error
		allowedCiphers <- tlsCipher
	}
}
