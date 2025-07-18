package dnsScan

import "net"

func IsIpv4(ip string) bool {
	// return strings.Count(ip, ":") < 2 // Explanation why this is accurate at https://stackoverflow.com/a/48519490
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}
