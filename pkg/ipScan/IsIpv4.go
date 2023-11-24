package ipScan

import "strings"

func IsIpv4(ip string) bool {
	return strings.Count(ip, ":") < 2 // Explanation why this is accurate at https://stackoverflow.com/a/48519490
}
