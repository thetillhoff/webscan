package ipScan

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func IsIpBlacklisted(ip string, verbose bool) ([]string, error) {
	var (
		err error

		resolver *net.Resolver

		reverseIp string = "" // includes trailing '.'

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

	for _, snippet := range strings.Split(ip, ".") {
		reverseIp = snippet + "." + reverseIp
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

		if verbose {
			fmt.Println("Checking ip blacklisting via", reverseIp+blacklist)
		}

		response, err = resolver.LookupIP(context.Background(), "ip4", reverseIp+blacklist)
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
				fmt.Println("Couldn't check ip blacklisting because of error code", response) // Display error
			} else { // If response isn't in error range
				blacklistsWithMatches = append(blacklistsWithMatches, blacklist) // Add blacklist match
			}
		}
	}

	return blacklistsWithMatches, nil
}
