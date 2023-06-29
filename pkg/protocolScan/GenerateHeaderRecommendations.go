package protocolScan

import (
	"net/http"
	"strings"
)

func GenerateHeaderRecommendations(response *http.Response) []string {
	var (
		err      error
		messages = []string{}

		headerName  string
		headerValue string
	)

	// HSTS
	headerName = "Strict-Transport-Security"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		messages = append(messages, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#http-strict-transport-security")
	} else {
		err = validateHeaderStrictTransportSecurity(headerValue)
		if err != nil {
			messages = append(messages, "Recommended action for "+headerName+": "+err.Error())
		}
	}

	// Content-Security-Policy
	headerName = "Content-Security-Policy"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		messages = append(messages, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#content-security-policy")
	} else {
		messages = append(messages, headerName+" header: "+headerValue) // TODO instead of just printing the header value check it against the best practices described in the link above
	}

	// X-Frame-Options
	headerName = "X-Frame-Options"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		messages = append(messages, headerName+" header should be set to 'sameorigin' or 'deny' as described in: https://infosec.mozilla.org/guidelines/web_security#x-frame-options")
	} else if headerValue == "sameorigin" || headerValue == "deny" {
		// Config okay
	} else {
		messages = append(messages, headerName+" should be set to 'sameorigin' or 'deny', but got '"+headerValue+"'")
	}

	// X-Content-Type-Options
	headerName = "X-Content-Type-Options"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		messages = append(messages, headerName+" header should be set to 'nosniff' as described in: https://infosec.mozilla.org/guidelines/web_security#x-content-type-options")
	} else if headerValue == "nosniff" {
		// Perfectly configured
	} else {
		messages = append(messages, headerName+" should be set to 'nosniff', but got '"+headerValue+"'")
	}

	// Referrer
	headerName = "Referer"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		// messages = append(messages, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#referrer-policy")
		// Default is used, which is strict-origin-when-cross-origin and therefore okay
	} else if headerValue == "no-referrer" || headerValue == "same-origin" || headerValue == "strict-origin" || headerValue == "strict-origin-when-cross-origin" {
		// Config okay
	} else {
		messages = append(messages, headerName+" should be set to one of 'no-referrer', 'same-origin', 'strict-origin', 'strict-origin-when-cross-origin', but got '"+headerValue+"'")
	}

	// Cache-Control
	headerName = "Cache-Control"
	headerValue = response.Header.Get(headerName)
	if headerValue == "" {
		messages = append(messages, headerName+" header should be configured, for example as described here: https://medium.com/pixelpoint/best-practices-for-cache-control-settings-for-your-website-ff262b38c5a2")
	} else {
		messages = append(messages, headerName+" header: "+headerValue) // TODO instead of just printing the header value check it against the best practices described in the link above
	}

	return messages
}
