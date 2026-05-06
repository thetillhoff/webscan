package ipScan

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

var spamhausListingCodes = map[string]string{
	"127.0.0.2":  "SBL - direct spam source",
	"127.0.0.3":  "SBL/CSS - compromised/spamming host",
	"127.0.0.4":  "XBL/CBL - exploited host (bot/proxy/trojan)",
	"127.0.0.5":  "XBL - reserved",
	"127.0.0.6":  "XBL - reserved",
	"127.0.0.7":  "XBL - reserved",
	"127.0.0.8":  "SBL - reserved",
	"127.0.0.9":  "SBL/DROP - hijacked netblock (do not route)",
	"127.0.0.10": "PBL - ISP maintained (end-user IP range)",
	"127.0.0.11": "PBL - Spamhaus maintained (end-user IP range)",
}

var spamhausErrorCodes = map[string]string{
	"127.255.255.252": "typing error in DNSBL zone name",
	"127.255.255.254": "query sent via public/open resolver (not allowed on free tier)",
	"127.255.255.255": "excessive number of queries (rate limited)",
}

func IsIPBlacklisted(ip string, timeout time.Duration) ([]string, error) {
	var (
		err error

		resolver *net.Resolver

		searchPrefix = ""
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

	if types.IsIPv4(ip) {
		network = "ip4"
		for _, snippet := range strings.Split(ip, ".") {
			searchPrefix = snippet + "." + searchPrefix
		}
	} else {
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
				d := net.Dialer{Timeout: timeout}
				return d.DialContext(ctx, network, net.JoinHostPort(blacklistNameserver, "53"))
			},
		}

		slog.Debug("ipScan: Checking for ip blacklisting", "blacklist", searchPrefix+blacklist)

		response, err = resolver.LookupIP(context.Background(), network, searchPrefix+blacklist)
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
			continue
		} else if err != nil {
			return blacklistsWithMatches, err
		}

		for _, responseIP := range response {
			code := responseIP.String()

			if errMsg, isError := spamhausErrorCodes[code]; isError {
				slog.Warn("Blacklist check returned error code", "ip", ip, "blacklist", blacklist, "code", code, "meaning", errMsg)
				continue
			}

			if listing, isListing := spamhausListingCodes[code]; isListing {
				blacklistsWithMatches = append(blacklistsWithMatches, fmt.Sprintf("%s (%s)", blacklist, listing))
			} else {
				blacklistsWithMatches = append(blacklistsWithMatches, fmt.Sprintf("%s (unknown code %s)", blacklist, code))
			}
		}
	}

	slog.Debug("ipScan: Checking for ip blacklisting completed")

	return blacklistsWithMatches, nil
}
