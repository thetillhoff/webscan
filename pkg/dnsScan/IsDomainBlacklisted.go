package dnsScan

import (
	"context"
	"log/slog"
	"net"
)

func IsDomainBlacklisted(domain string, resolver *net.Resolver) ([]string, error) {
	var (
		err        error
		blacklists = []string{
			"zen.spamhaus.org",
		}

		blacklistsWithMatches = []string{}
	)

	slog.Debug("dnsScan: Checking for domain blacklisting started")

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

	slog.Debug("dnsScan: Checking for domain blacklisting completed")

	return blacklistsWithMatches, nil
}
