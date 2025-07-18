package tlsScan

import (
	"crypto/tls"
	"slices"
	"strings"
)

type certInfo struct {
	names   []string
	issuers []string
}

type Result struct {
	tlsScanResultPerIp map[string]TlsScanResult
	// sharedCipherRules  map[string][]string
}

type TlsScanResult struct {
	certInfos          []certInfo
	tlsErr             error
	enabledTlsVersions []uint16
	enabledTlsCiphers  []tls.CipherSuite

	cipherRulesEvaluationResult map[string][]string
}

// Returns list of cert names for an ip
func (r *Result) ListCertNamesForIp(ip string) []string {

	certNames := []string{}
	if len(r.tlsScanResultPerIp[ip].certInfos) > 0 {
		certNames = append(certNames, r.tlsScanResultPerIp[ip].certInfos[0].names...) // Only the first one, since others are for the certificate chain
	}
	slices.Sort(certNames)
	return slices.Compact(certNames)
}

// Returns list of all cert names for all ips
func (r *Result) ListAllCertNames() []string {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	var certNames []string
	for ip := range r.tlsScanResultPerIp {
		ipCertNames := r.ListCertNamesForIp(ip)
		certNames = append(certNames, ipCertNames...)
	}

	return certNames
}

// Returns list of shared cert names for all ips
func (r *Result) ListSharedCertNames() []string {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get cert names from first IP as initial set
	var certNames []string
	for ip := range r.tlsScanResultPerIp {
		certNames = r.ListCertNamesForIp(ip)
		break
	}

	// For each remaining IP, keep only cert names that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipCertNames := r.ListCertNamesForIp(ip)
		certNames = slices.DeleteFunc(certNames, func(name string) bool {
			return !slices.Contains(ipCertNames, name)
		})

		if len(certNames) == 0 {
			break // Early exit if no shared names remain
		}
	}

	return certNames
}

// Returns list of non-shared cert names for an ip
func (r *Result) ListNonSharedCertNamesForIp(ip string) []string {
	certNamesForIp := r.ListCertNamesForIp(ip)
	nonSharedCertNames := []string{}
	for _, certName := range certNamesForIp {
		if !slices.Contains(r.ListSharedCertNames(), certName) {
			nonSharedCertNames = append(nonSharedCertNames, certName)
		}
	}
	return nonSharedCertNames
}

// Returns list of cert issuers for an ip
func (r *Result) ListCertIssuersForIp(ip string) []string {
	certIssuers := []string{}
	if len(r.tlsScanResultPerIp[ip].certInfos) > 0 {
		certIssuers = append(certIssuers, r.tlsScanResultPerIp[ip].certInfos[0].issuers...) // Only the first one, since others are for the certificate chain
	}
	slices.Sort(certIssuers)
	return slices.Compact(certIssuers)
}

// Returns list of shared cert issuers for all ips
func (r *Result) ListSharedCertIssuers() []string {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get cert issuers from first IP as initial set
	var certIssuers []string
	for ip := range r.tlsScanResultPerIp {
		certIssuers = r.ListCertIssuersForIp(ip)
		break
	}

	// For each remaining IP, keep only cert issuers that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipCertIssuers := r.ListCertIssuersForIp(ip)
		certIssuers = slices.DeleteFunc(certIssuers, func(issuer string) bool {
			return !slices.Contains(ipCertIssuers, issuer)
		})

		if len(certIssuers) == 0 {
			break // Early exit if no shared issuers remain
		}
	}

	return certIssuers
}

// Returns list of non-shared cert issuers for an ip
func (r *Result) ListNonSharedCertIssuersForIp(ip string) []string {
	certIssuersForIp := r.ListCertIssuersForIp(ip)
	nonSharedCertIssuers := []string{}
	for _, certIssuer := range certIssuersForIp {
		if !slices.Contains(r.ListSharedCertIssuers(), certIssuer) {
			nonSharedCertIssuers = append(nonSharedCertIssuers, certIssuer)
		}
	}
	return nonSharedCertIssuers
}

// Returns tls error for an ip
func (r *Result) GetTlsErrForIp(ip string) error {
	return r.tlsScanResultPerIp[ip].tlsErr
}

// Returns list of shared tls errors for all ips
func (r *Result) ListSharedTlsErr() []error {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get tls error from first IP as initial set
	var tlsErrs []error
	for ip := range r.tlsScanResultPerIp {
		tlsErr := r.GetTlsErrForIp(ip)
		if tlsErr != nil {
			tlsErrs = append(tlsErrs, tlsErr)
		}
		break
	}

	// For each remaining IP, keep only cert names that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipTlsErr := r.GetTlsErrForIp(ip)
		tlsErrs = slices.DeleteFunc(tlsErrs, func(err error) bool {
			return ipTlsErr != err
		})

		if len(tlsErrs) == 0 {
			break // Early exit if no shared errors remain
		}
	}

	return tlsErrs
}

// Returns list of non-shared tls errors for an ip
func (r *Result) ListNonSharedTlsErrForIp(ip string) []error {
	tlsErrsForIp := r.GetTlsErrForIp(ip)
	nonSharedTlsErrs := []error{}
	if tlsErrsForIp != nil {
		if !slices.ContainsFunc(r.ListSharedTlsErr(), func(err error) bool {
			return tlsErrsForIp == err
		}) {
			nonSharedTlsErrs = append(nonSharedTlsErrs, tlsErrsForIp)
		}
	}
	return nonSharedTlsErrs
}

// Returns list of tls versions for an ip
func (r *Result) ListTlsVersionsForIp(ip string) []uint16 {
	tlsVersions := []uint16{}
	if len(r.tlsScanResultPerIp[ip].enabledTlsVersions) > 0 {
		tlsVersions = append(tlsVersions, r.tlsScanResultPerIp[ip].enabledTlsVersions...)
	}
	slices.Sort(tlsVersions)
	return slices.Compact(tlsVersions)
}

// Returns list of shared tls versions for all ips
func (r *Result) ListSharedTlsVersions() []uint16 {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get tls versions from first IP as initial set
	var tlsVersions []uint16
	for ip := range r.tlsScanResultPerIp {
		tlsVersions = r.ListTlsVersionsForIp(ip)
		break
	}

	// For each remaining IP, keep only tls versions that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipTlsVersions := r.ListTlsVersionsForIp(ip)
		tlsVersions = slices.DeleteFunc(tlsVersions, func(version uint16) bool {
			return !slices.Contains(ipTlsVersions, version)
		})

		if len(tlsVersions) == 0 {
			break // Early exit if no shared versions remain
		}
	}

	return tlsVersions
}

// Returns list of non-shared tls versions for an ip
func (r *Result) ListNonSharedTlsVersionsForIp(ip string) []uint16 {
	tlsVersionsForIp := r.ListTlsVersionsForIp(ip)
	nonSharedTlsVersions := []uint16{}
	for _, tlsVersion := range tlsVersionsForIp {
		if !slices.Contains(r.ListSharedTlsVersions(), tlsVersion) {
			nonSharedTlsVersions = append(nonSharedTlsVersions, tlsVersion)
		}
	}
	return nonSharedTlsVersions
}

// Returns list of tls ciphers for an ip
func (r *Result) ListTlsCiphersForIp(ip string) []tls.CipherSuite {
	tlsCiphers := []tls.CipherSuite{}
	if len(r.tlsScanResultPerIp[ip].enabledTlsCiphers) > 0 {
		tlsCiphers = append(tlsCiphers, r.tlsScanResultPerIp[ip].enabledTlsCiphers...)
	}
	slices.SortFunc(tlsCiphers, func(a, b tls.CipherSuite) int {
		return strings.Compare(a.Name, b.Name)
	})
	return slices.CompactFunc(tlsCiphers, func(a, b tls.CipherSuite) bool {
		return a.ID == b.ID
	})
}

// Returns list of shared tls ciphers for all ips
func (r *Result) ListSharedTlsCiphers() []tls.CipherSuite {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get tls versions from first IP as initial set
	var tlsCiphers []tls.CipherSuite
	for ip := range r.tlsScanResultPerIp {
		tlsCiphers = r.ListTlsCiphersForIp(ip)
		break
	}

	// For each remaining IP, keep only tls versions that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipTlsCiphers := r.ListTlsCiphersForIp(ip)
		tlsCiphers = slices.DeleteFunc(tlsCiphers, func(cipher tls.CipherSuite) bool {
			return !slices.ContainsFunc(ipTlsCiphers, func(c tls.CipherSuite) bool {
				return c.ID == cipher.ID
			})
		})

		if len(tlsCiphers) == 0 {
			break // Early exit if no shared ciphers remain
		}
	}

	return tlsCiphers
}

// Returns list of non-shared tls ciphers for an ip
func (r *Result) ListNonSharedTlsCiphersForIp(ip string) []tls.CipherSuite {
	tlsCiphersForIp := r.ListTlsCiphersForIp(ip)
	nonSharedTlsCiphers := []tls.CipherSuite{}
	for _, tlsCipher := range tlsCiphersForIp {
		if !slices.ContainsFunc(r.ListSharedTlsCiphers(), func(c tls.CipherSuite) bool {
			return c.ID == tlsCipher.ID
		}) {
			nonSharedTlsCiphers = append(nonSharedTlsCiphers, tlsCipher)
		}
	}
	return nonSharedTlsCiphers
}

// Returns list of cipher rules for an ip
func (r *Result) ListCipherRulesForIp(ip string) map[string][]string {
	cipherRuleEvaluationResults := map[string][]string{}
	if len(r.tlsScanResultPerIp[ip].cipherRulesEvaluationResult) > 0 {
		for rule := range r.tlsScanResultPerIp[ip].cipherRulesEvaluationResult {
			cipherRuleEvaluationResults[rule] = r.tlsScanResultPerIp[ip].cipherRulesEvaluationResult[rule]
		}
	}
	return cipherRuleEvaluationResults
}

// Returns list of shared cipher rules for all ips
func (r *Result) ListSharedCipherRules() map[string][]string {
	if len(r.tlsScanResultPerIp) == 0 {
		return nil
	}

	// Get tls versions from first IP as initial set
	var cipherRuleEvaluationResults map[string][]string
	for ip := range r.tlsScanResultPerIp {
		cipherRuleEvaluationResults = r.ListCipherRulesForIp(ip)
		break
	}

	// For each remaining IP, keep only tls versions that exist for all IPs
	for ip := range r.tlsScanResultPerIp {
		ipCipherRuleEvaluationResults := r.ListCipherRulesForIp(ip)
		for rule := range cipherRuleEvaluationResults {
			cipherRuleEvaluationResults[rule] = slices.DeleteFunc(cipherRuleEvaluationResults[rule], func(cipher string) bool {
				return !slices.Contains(ipCipherRuleEvaluationResults[rule], cipher)
			})
		}

		if len(cipherRuleEvaluationResults) == 0 {
			break // Early exit if no shared rules remain
		}
	}

	return cipherRuleEvaluationResults
}

// Returns list of non-shared cipher rules for an ip
func (r *Result) ListNonSharedCipherRulesForIp(ip string) map[string][]string {
	cipherRuleEvaluationResultsForIp := r.ListCipherRulesForIp(ip)
	sharedRules := r.ListSharedCipherRules()
	nonSharedCipherRuleEvaluationResults := map[string][]string{}

	for rule, ciphers := range cipherRuleEvaluationResultsForIp {
		// Check if rule exists in shared rules and has different ciphers
		if sharedCiphers, exists := sharedRules[rule]; !exists || !slices.Equal(ciphers, sharedCiphers) {
			nonSharedCipherRuleEvaluationResults[rule] = ciphers
		}
	}
	return nonSharedCipherRuleEvaluationResults
}
