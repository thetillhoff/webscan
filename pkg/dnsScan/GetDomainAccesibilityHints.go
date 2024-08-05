package dnsScan

import (
	"log/slog"
	"strings"
)

func GetDomainAccessibilityHints(url string) []string {
	var (
		hints = []string{}
	)

	slog.Debug("dnsScan: Getting domain accessibility hints started")

	domains := strings.Split(url, ".")

	if len(url) > 20 { // Url shouldn't be too long
		hints = append(hints, "Hint: Url `"+url+"` has quite many characters.")
	}

	for _, domain := range domains {

		words := strings.Split(domain, "-")
		if len(words) > 3 { // Single part of a domain should not have too many words
			hints = append(hints, "Hint: Domain `"+domain+"` has quite many words.")
		}

		for _, word := range words {
			if len(word) > 12 { // Words within domain part should not be too long
				hints = append(hints, "Hint: Domain `"+domain+"` contains a quite long word.")
			}
		}
	}

	slog.Debug("dnsScan: Getting domain accessibility hints started")

	return hints
}
