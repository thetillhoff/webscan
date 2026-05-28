package tlsScan

import (
	"crypto/tls"
	"log/slog"
)

// verify ciphers against best practices
// The result is a map of rule description to list of ciphers that match the rule
func evaluateTLSCipherRules(tlsCipherSuites []tls.CipherSuite) map[string][]string {
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

				if _, ok := ruleMatches[rule.title]; !ok { // If map entry doesn't exist
					ruleMatches[rule.title] = []string{} // Initialize map entry
				}

				ruleMatches[rule.title] = append(ruleMatches[rule.title], tlsCipherSuite.Name) // Add cipherSuite name to list
			}
		}

	}

	// Post-processing: remove from the Golang entry any ciphers already matched by another rule
	const golangTitle = "Ciphers deemed insecure by Golang"
	if golangMatches, ok := ruleMatches[golangTitle]; ok {
		// Collect all cipher names matched by non-Golang rules
		otherMatched := map[string]bool{}
		for title, names := range ruleMatches {
			if title != golangTitle {
				for _, name := range names {
					otherMatched[name] = true
				}
			}
		}
		// Filter the Golang entry
		filtered := golangMatches[:0]
		for _, name := range golangMatches {
			if !otherMatched[name] {
				filtered = append(filtered, name)
			}
		}
		if len(filtered) == 0 {
			delete(ruleMatches, golangTitle)
		} else {
			ruleMatches[golangTitle] = filtered
		}
	}

	slog.Debug("tlsScan: Evaluating cipher rules completed")

	return ruleMatches
}
