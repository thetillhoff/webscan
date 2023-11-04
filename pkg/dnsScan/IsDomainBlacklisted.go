package dnsScan

import (
	"context"
	"net"
)

func (engine Engine) IsDomainBlacklisted(domain string, resolver *net.Resolver) ([]string, error) {
	var (
		err        error
		blacklists = []string{
			"zen.spamhaus.org",
		}

		blacklistsWithMatches = []string{}
	)

	for _, blacklist := range blacklists {
		_, err = resolver.LookupIP(context.Background(), "ip4", domain+"."+blacklist)
		if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
			// No A record available -> Not blacklisted
		} else if err != nil { // Unknown error occurred
			return blacklistsWithMatches, err
		} else { // Record was found
			blacklistsWithMatches = append(blacklistsWithMatches, blacklist)
		}
	}

	return blacklistsWithMatches, nil
}
