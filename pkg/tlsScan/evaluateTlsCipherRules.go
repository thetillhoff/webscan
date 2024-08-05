package tlsScan

import (
	"crypto/tls"
	"log/slog"
)

// verify ciphers against best practices
func evaluateTlsCipherRules(tlsCipherSuites []tls.CipherSuite) map[string][]string {
	var (
		rules       = getRules()
		ruleMatches = map[string][]string{}
	)

	slog.Debug("tlsScan: Evaluating cipher rules started")

	// Idea is to iterate over rules, and for each rule iterate over ciphers
	// If one of the ciphers matches the rule, add the rule to matchedRules and list of

	for _, rule := range rules {

		// Verify ciphers (https://ciphersuite.info/cs/?tls=tls12&singlepage=true has some nice hints on the reasons behind deeming a cipher insecure)
		for _, tlsCipherSuite := range tlsCipherSuites {

			if rule.matchFunc(tlsCipherSuite) {

				if _, ok := ruleMatches[rule.description]; !ok { // If map entry doesn't exist
					ruleMatches[rule.description] = []string{} // Initialize map entry
				}

				ruleMatches[rule.description] = append(ruleMatches[rule.description], tlsCipherSuite.Name) // Add cipherSuite name to list
			}
		}

	}

	slog.Debug("tlsScan: Evaluating cipher rules completed")

	return ruleMatches
}
