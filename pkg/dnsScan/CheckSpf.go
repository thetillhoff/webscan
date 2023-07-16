package dnsScan

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

// TODO
// Add `--evaluate-mail-server <mail-server>` that will be checked against the verified spf record

func (engine Engine) CheckSpf() string {
	var (
		spfRecord string = ""
		word      string

		existingExp      bool = false
		existingRedirect bool = false
	)

	for _, txtRecord := range engine.TXTRecords {
		if txtRecord == "v=spf1" || strings.HasPrefix(txtRecord, "v=spf1 ") {
			if spfRecord == "" { // Check if there was a spf record detected before
				spfRecord = strings.ToLower(spfRecord) // spf records are case insensitive
			} else {
				return "Hint: Multiple SPF records detected."
			}
		}
	}

	if spfRecord == "" { // Check if there is at least one
		return "Hint: No SPF record detected."
	}

	spfRecord = strings.TrimPrefix(spfRecord, "v=spf1") // Remove spf prefix, since it's always the same
	words := strings.Split(spfRecord, " ")

	fmt.Println(words)

	qualifiers := []string{"", "+", "-", "?", "~"}
	mechanisms := []string{"all", "include", "a", "mx", "ptr", "ip4", "ip6", "exists"}

	for len(words) > 0 {
		word, words = words[0], words[1:] // Get word for processing and remove it from list

		if listContains(generateCartesianProduct(qualifiers, []string{"all"}), word) {
			for _, remainingWord := range words { // Mechanisms listed after "all" MUST be ignored.
				for _, prefix := range generateCartesianProduct(qualifiers, mechanisms) {
					if strings.HasPrefix(remainingWord, prefix) {
						return "Hint: SPF mechanisms like `" + remainingWord + "` that are listed after `all` are ignored."
					}
				}
			}

			// Any "redirect" modifier (Section 6.1) MUST be ignored when there is an "all" mechanism in the record, regardless of the relative ordering of the terms.
			for _, remainingWord := range words {
				for _, qualifier := range qualifiers { // Remove qualifier prefix if exists
					remainingWord = strings.TrimPrefix(remainingWord, qualifier)
				}
				if strings.HasPrefix(remainingWord, "redirect=") {
					return "Hint: SPF modifiers like `redirect` are ignored when placed after the `all` mechanism."
				}
			}
		} else if listContains(generateCartesianProduct(qualifiers, []string{"include"}), word) {
			// TODO
			// "include:<domain-spec>"
			// Include will check whether the spf record of the specified domain matches. if it does, this include matches and processing ends.
			// It might make sense to verify SPF records in an recursive way because of this.
		} else if listContains(generateCartesianProduct(qualifiers, []string{"a"}), word) {
			// a[:<domain-spec>][<dual-cidr-length>]
			// the "a" mechanism also matches AAAA records.

			if len(word) > 1 { // Word contains at least <domain-spec>
				return "Hint: SPF record contains a rather complex setup around the `a` mechanism."

				// word = strings.TrimPrefix(word, "a:")

				// TODO remove <domain-spec> from word, and only continue if the remaining string is not empty

				// if strings.Contains(word, "/") { // Contains both <ipv4-cidr-length> and <ipv6-cidr-length>
				// 	dualCidrs := strings.SplitN(word, "/", 2)
				// 	ipv4Cidr, ipv6Cidr := dualCidrs[0], dualCidrs[1]
				// 	_, _, err = net.ParseCIDR(ipv4Cidr)
				// 	if err != nil || strings.Count(ipv4Cidr, ".") != 3 { // Check for parsing error and crude check for ipv4 vs ipv6
				// 		log.Fatalln("Invalid IPv4 CIDR range in SPF record:", err)
				// 	}

				// 	_, _, err = net.ParseCIDR(ipv6Cidr)
				// 	if err != nil || strings.Count(ipv4Cidr, ".") != 3 { // Check for parsing error and crude check for ipv4 vs ipv6
				// 		log.Fatalln("Invalid IPv6 CIDR range in SPF record:", err)
				// 	}
				// } else { // Contains only <ipv4-cidr-length>
				// 	_, _, err := net.ParseCIDR(word)
				// 	if err != nil || strings.Count(word, ".") == 3 { // Check for parsing error and crude check for ipv4 vs ipv6
				// 		log.Fatalln("Invalid IPv4 CIDR range in SPF record:", err)
				// 	}
				// }
			}
		} else if listContains(generateCartesianProduct(qualifiers, []string{"mx"}), word) {
			// TODO
			// mx[:<domain-spec>][<dual-cidr-length>]
			// the "mx" mechanism also matches AAAA records.
		} else if listContains(generateCartesianProduct(qualifiers, []string{"ptr"}), word) {
			// The "ptr" mechanism SHOULD NOT be used because it is slow and inefficient.
			return "Hint: PTR records should not be used inSPF records as they are slow and inefficient."

		} else if listContains(generateCartesianProduct(qualifiers, []string{"ip4"}), word) {

			for _, qualifier := range qualifiers { // Remove qualifier prefix if exists
				word = strings.TrimPrefix(word, qualifier)
			}

			// "ip4:<ip4-network>[<ip4-cidr-length>]
			word = strings.TrimPrefix(word, "ipv4:") // Remove mechanism prefix

			if strings.Contains(word, "/") { // Contains both <ipv4-network> and <ipv4-cidr-length>
				// ip4-network      = qnum "." qnum "." qnum "." qnum
				// qnum             = 0 - 255
				// ip4-cidr-length  = "/" ("0" / %x31-39 0*1DIGIT) ; value range 0-32
				_, _, err := net.ParseCIDR(word)
				if err != nil || !IsIpv4(word) { // Check if ip is really ipv4, not ipv6
					log.Fatalln("Invalid IPv4 address / cidr-range in SPF record:", word)
				}
				// It is not permitted to omit parts of the IP address instead of using CIDR notations. That is, use 192.0.2.0/24 instead of 192.0.2.

			} else { // Contains only <ipv4-network>
				// ip4-network      = qnum "." qnum "." qnum "." qnum
				// qnum             = 0 - 255
				// If ip4-cidr-length is omitted, it is taken to be "/32".
				parsedIp := net.ParseIP(word)
				if parsedIp == nil || !IsIpv4(word) { // Check if ip is really ipv4, not ipv6
					log.Fatalln("Invalid IPv4 address in SPF record:", word)
				}
			}

		} else if listContains(generateCartesianProduct(qualifiers, []string{"ip6"}), word) {

			for _, qualifier := range qualifiers { // Remove qualifier prefix if exists
				word = strings.TrimPrefix(word, qualifier)
			}

			// "ip6:<ip6-network>[<ip6-cidr-length>]
			word = strings.TrimPrefix(word, "ipv6:") // Remove mechanism prefix

			if strings.Contains(word, "/") { // Contains both <ipv6-network> and <ipv6-cidr-length>
				// ip6-network      = <as per Section 2.2 of [RFC4291]>
				// ip6-cidr-length  = "/" ("0" / %x31-39 0*2DIGIT) ; value range 0-128
				_, _, err := net.ParseCIDR(word)
				if err != nil || IsIpv4(word) {
					log.Fatalln("Invalid IPv6 address / cidr-range in SPF record:", word)
				}
				// It is not permitted to omit parts of the IP address instead of using CIDR notations. That is, use 192.0.2.0/24 instead of 192.0.2.

			} else { // Contains only <ipv6-network>
				// ip6-network      = <as per Section 2.2 of [RFC4291]>
				// If ip6-cidr-length is omitted, it is taken to be "/128".
				parsedIp := net.ParseIP(word)
				if parsedIp == nil || IsIpv4(word) { // Check if ip is really ipv6, not ipv4
					log.Fatalln("Invalid IPv6 address in SPF record:", word)
				}
			}

		} else if listContains(generateCartesianProduct(qualifiers, []string{"exists"}), word) {
			// TODO
			// "exists:<domain-spec>"
			// The resulting domain name is used for a DNS A RR lookup (even when the connection type is IPv6).  If any A record is returned, this mechanism matches.
		} else if strings.Contains(word, "=") { // Modifiers always have an "=" separating the name and the value.
			if strings.HasPrefix(word, "redirect=") {
				// "redirect=<domain-spec>"
				word = strings.TrimPrefix(word, "redirect=") // Remove modifier prefix

				// TODO verify <domain-spec>

				// if no SPF record is found, or if the <target-name> is malformed, the result is a "permerror" rather than "none".
				// Note that the newly queried domain can itself specify redirect processing.
				// any "redirect" modifier SHOULD appear as the very last term in a record.  Any "redirect" modifier MUST be ignored if there is an "all" mechanism anywhere in the record.

				// Modifiers ("redirect" and "exp") MUST NOT appear in a record more than once each.
				if existingRedirect {
					return "Hint: SPF modifier `redirect` must not appear more than once."
				} else {
					existingRedirect = true
				}
				if strings.HasPrefix(word, "exp=") {
					// "exp=<domain-spec>" describes _why_ something failed. There are defaults in place however
					// word = strings.TrimPrefix(word, "exp=") // Remove modifier prefix

					// TODO verify <domain-spec>

					// The DNS TXT RRset for the <target-name> is fetched. If there are any DNS processing errors (any RCODE other than 0), or if no records are returned, or if more than one record is returned, or if there are syntax errors in the explanation string, then proceed as if no "exp" modifier was given.
					// fetched TXT record's strings are concatenated with no spaces which is macro-expanded. This final result is the explanation string.
					// Since the explanation string is intended for an SMTP response and Section 2.4 of [RFC5321] says that responses are in [US-ASCII], the explanation string MUST be limited to [US-ASCII].
					// Example usage: "v=spf1 mx -all exp=explain._spf.%{d}""
					//   Here are some examples of possible explanation TXT records at explain._spf.example.com:
					//   "Mail from example.com should only be sent by its own servers."
					//   "%{i} is not one of %{d}'s designated mail servers."
					// During recursion into an "include" mechanism, an "exp" modifier from the <target-name> MUST NOT be used.
					// In contrast, when executing a "redirect" modifier, an "exp" modifier from the original domain MUST NOT be used.  This is because "include" is meant to cross administrative boundaries and the explanation provided should be the one from the receiving ADMD, while "redirect" is meant to operate as a tool to consolidate policy records within an ADMD so the redirected explanation is the one that ought to have priority.

					// Modifiers ("redirect" and "exp") MUST NOT appear in a record more than once each.
					if existingExp {
						return "Hint: SPF modifier `exp` must not appear more than once."
					} else {
						existingExp = true
					}
				} else { // Unkown modifier
					// Unrecognized modifiers MUST be ignored
					// unknown-modifier = "<name>=<macro-string>" ; where name is not "redirect" or "exp"
					// name             = ALPHA *( ALPHA / DIGIT / "-" / "_" / "." )

					modifier := strings.SplitN(word, "=", 2)
					name, macroString := modifier[0], modifier[1]

					matched, err := regexp.MatchString(`^[[:alpha:]]([[:alnum:]]|-|_|.)*$`, name)
					if err != nil {
						log.Fatalln(err)
					}

					if !matched {
						return "Hint: Unkown SPF modifier name `" + name + "` is invalid."
					} else {
						return "Hint: Unkown SPF modifier `" + name + "=" + macroString + "`." // TODO validate macroString instead
					}
				}
			} else { // Invalid setting
				return "Hint: SPF configuration invalid, due to unkown setting `" + word + "`."
			}
		}
	}

	// dual-cidr-length = [ ip4-cidr-length ] [ "/" ip6-cidr-length ]

	// Modifiers

	// The following terms cause DNS queries:
	// "include", "a", "mx", "ptr", and "exists" mechanisms, and the "redirect" modifier.
	// SPF implementations MUST limit the total number of those terms to 10 during SPF evaluation, to avoid unreasonable load on the DNS.
	// When evaluating the "mx" mechanism, the number of "MX" resource records queried is included in the overall limit of 10 mechanisms/modifiers that cause DNS lookups as described above.
	// In addition to that limit, the evaluation of each "MX" record MUST NOT result in querying more than 10 address records -- either "A" or "AAAA" resource records.

	// For several mechanisms, the <domain-spec> is optional. If it is not provided, the <domain> is used as the <target-name>.

	// If none of the mechanisms match and there is no "redirect" modifier, then the check_host() returns a result of "neutral", just as if "?all" were specified as the last directive.  If there is a "redirect" modifier, check_host() proceeds as defined in Section 6.1.
	// It is better to use either a "redirect" modifier or an "all" mechanism to explicitly terminate processing.  Although there is an implicit "?all" at the end of every record that is not explicitly terminated, it aids debugging efforts when it is explicitly provided.

	return "" // No issues found
}
