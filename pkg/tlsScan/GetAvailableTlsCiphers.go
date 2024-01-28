package tlsScan

import (
	"crypto/tls"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// Checks whether the specified <url> uses a secure configuration
// <url> has to be in the format of <ip>:<port>
func GetAvailableTlsCiphers(url string) []TlsCipher {
	var (
		tlsVersions = []uint16{
			tls.VersionTLS13,
			tls.VersionTLS12,
			tls.VersionTLS11,
			tls.VersionTLS10,
		}
		ciphers = []uint16{}

		allowedTlsCiphersChan = make(chan TlsCipher, len(tlsVersions)*(len(tls.CipherSuites())+len(tls.InsecureCipherSuites())))
	)

	urlHost := strings.SplitN(url, "/", 2)[0] // Cleaning url (removing Path)

	for _, cipher := range tls.CipherSuites() { // Add the secure ciphers
		ciphers = append(ciphers, cipher.ID)
	}
	for _, cipher := range tls.InsecureCipherSuites() { // Add the insecure ciphers
		ciphers = append(ciphers, cipher.ID)
	}

	for _, tlsVersion := range tlsVersions { // For each tls version
		for _, cipher := range ciphers { // For each cipher
			wg.Add(1)                                                                                               // Wait for one more goroutine to finish
			go isCipherAvailable(urlHost, TlsCipher{TlsVersion: tlsVersion, Cipher: cipher}, allowedTlsCiphersChan) // Start goroutine that checks if tlsVersion and cipher combination are allowed
		}
	}

	wg.Wait() // Wait until all goroutines are finished

	close(allowedTlsCiphersChan) // Make sure channel is closed when goroutines are finished

	allowedTlsCiphers := []TlsCipher{}

	for allowedTlsCipher := range allowedTlsCiphersChan {
		allowedTlsCiphers = append(allowedTlsCiphers, allowedTlsCipher)
	}

	return allowedTlsCiphers
}
