package ipScan

import (
	"context"
	"log/slog"
	"math/rand"
	"net"
	"net/netip"
	"strings"
	"time"
)

func IsIpBlacklisted(ip string) ([]string, error) {
	var (
		err error

		resolver *net.Resolver

		searchPrefix = "" // includes trailing '.'
		network      string

		blacklistWithNameservers = map[string][]string{
			"zen.spamhaus.org": {
				"a.gns.spamhaus.org",
				"b.gns.spamhaus.org",
				"c.gns.spamhaus.org",
				"d.gns.spamhaus.org",
				"e.gns.spamhaus.org",
			},
		}

		response []net.IP

		blacklistsWithMatches = []string{}
	)

	slog.Debug("ipScan: Checking for ip blacklisting started")

	if IsIpv4(ip) { // If ip is ipv4
		network = "ip4"
		for _, snippet := range strings.Split(ip, ".") {
			searchPrefix = snippet + "." + searchPrefix
		}
	} else { // If ip is ipv6
		network = "ip6"

		addr, _ := netip.ParseAddr(ip)
		ip = addr.StringExpanded()
		ip = strings.ReplaceAll(ip, ":", "")
		for _, snippet := range strings.Split(ip, "") {
			searchPrefix = snippet + "." + searchPrefix
		}
	}

	for blacklist, blacklistNameservers := range blacklistWithNameservers {
		blacklistNameserver := blacklistNameservers[rand.Intn(len(blacklistNameservers))]

		resolver = &net.Resolver{
			PreferGo:     false,
			StrictErrors: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
				return d.DialContext(ctx, network, net.JoinHostPort(blacklistNameserver, "53"))
			},
		}

		slog.Debug("ipScan: Checking for ip blacklisting", "blacklist", searchPrefix+blacklist)

		response, err = resolver.LookupIP(context.Background(), network, searchPrefix+blacklist)
		if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
			// No A record available -> Not blacklisted, so continue
		} else if err != nil { // Unknown error occurred
			return blacklistsWithMatches, err
		} else { // Record was found

			_, ipv4Net, err := net.ParseCIDR("127.255.255.0/24") // Parse CIDR of error range of spamhaus
			if err != nil {
				return blacklistsWithMatches, err
			}

			if len(response) == 1 && ipv4Net.Contains(response[0]) { // If response is in error range
				slog.Warn("Couldn't check ip blacklisting because of error code", "ip", ip, "response", response)
			} else { // If response isn't in error range
				blacklistsWithMatches = append(blacklistsWithMatches, blacklist) // Add blacklist match
			}
		}
	}

	slog.Debug("ipScan: Checking for ip blacklisting completed")

	return blacklistsWithMatches, nil
}
