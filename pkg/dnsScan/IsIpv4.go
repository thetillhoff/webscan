package dnsScan

import "strings"

func IsIpv4(ip string) bool {
	// TODO find better way to do this - is this even needed?
	return strings.Count(ip, ":") < 2 // Explanation why this is accurate at https://stackoverflow.com/a/48519490
}
