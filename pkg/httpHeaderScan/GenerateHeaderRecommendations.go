package httpHeaderScan

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func GenerateHeaderRecommendations(response *http.Response, schema types.Schema) []HeaderEntry {
	var entries []HeaderEntry

	slog.Debug("httpHeaderScan: Generating header recommendations started")

	// Server — informational only
	if v := strings.ToLower(response.Header.Get("Server")); v != "" {
		entries = append(entries, HeaderEntry{Name: "Server", Value: v})
	}

	// Strict-Transport-Security (HTTPS only)
	if schema == types.HTTPS {
		v := response.Header.Get("Strict-Transport-Security")
		if v == "" {
			entries = append(entries, HeaderEntry{
				Name:           "Strict-Transport-Security",
				Recommendation: "Should be implemented: https://infosec.mozilla.org/guidelines/web_security#http-strict-transport-security",
			})
		} else {
			entry := HeaderEntry{Name: "Strict-Transport-Security", Value: strings.ToLower(v)}
			if err := validateSTS(strings.ToLower(v)); err != nil {
				entry.Recommendation = err.Error()
			}
			entries = append(entries, entry)
		}
	}

	// Content-Security-Policy
	if v := response.Header.Get("Content-Security-Policy"); v == "" {
		entries = append(entries, HeaderEntry{
			Name:           "Content-Security-Policy",
			Recommendation: "Should be implemented: https://infosec.mozilla.org/guidelines/web_security#content-security-policy",
		})
	} else {
		var directives []string
		for _, d := range strings.Split(strings.ToLower(v), ";") {
			if d = strings.TrimSpace(d); d != "" {
				directives = append(directives, d)
			}
		}
		entries = append(entries, HeaderEntry{
			Name:  "Content-Security-Policy",
			Value: strings.Join(directives, ";\n"),
		})
	}

	// X-Frame-Options
	switch v := strings.ToLower(response.Header.Get("X-Frame-Options")); v {
	case "":
		entries = append(entries, HeaderEntry{
			Name:           "X-Frame-Options",
			Recommendation: "Should be set to 'sameorigin' or 'deny': https://infosec.mozilla.org/guidelines/web_security#x-frame-options",
		})
	case "sameorigin", "deny":
		// correctly configured
	default:
		entries = append(entries, HeaderEntry{
			Name:           "X-Frame-Options",
			Value:          v,
			Recommendation: "Should be 'sameorigin' or 'deny'",
		})
	}

	// X-Content-Type-Options
	switch v := strings.ToLower(response.Header.Get("X-Content-Type-Options")); v {
	case "":
		entries = append(entries, HeaderEntry{
			Name:           "X-Content-Type-Options",
			Recommendation: "Should be set to 'nosniff': https://infosec.mozilla.org/guidelines/web_security#x-content-type-options",
		})
	case "nosniff":
		// correctly configured
	default:
		entries = append(entries, HeaderEntry{
			Name:           "X-Content-Type-Options",
			Value:          v,
			Recommendation: "Should be 'nosniff'",
		})
	}

	// Referrer-Policy (note: "Referer" is a request header; "Referrer-Policy" is the response header)
	switch v := strings.ToLower(response.Header.Get("Referrer-Policy")); v {
	case "", "no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin",
		"same-origin", "strict-origin", "strict-origin-when-cross-origin":
		// not set uses browser default (strict-origin-when-cross-origin), all listed values are acceptable
	case "unsafe-url":
		entries = append(entries, HeaderEntry{
			Name:           "Referrer-Policy",
			Value:          v,
			Recommendation: "'unsafe-url' sends the full URL on all requests including cross-origin — consider 'strict-origin-when-cross-origin'",
		})
	default:
		entries = append(entries, HeaderEntry{
			Name:           "Referrer-Policy",
			Value:          v,
			Recommendation: "Unknown value; valid options: no-referrer, no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url",
		})
	}

	// Cache-Control
	if v := response.Header.Get("Cache-Control"); v == "" {
		entries = append(entries, HeaderEntry{
			Name:           "Cache-Control",
			Recommendation: "Should be configured: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control",
		})
	} else {
		entries = append(entries, HeaderEntry{Name: "Cache-Control", Value: v})
	}

	slog.Debug("httpHeaderScan: Generating header recommendations completed")

	return entries
}
