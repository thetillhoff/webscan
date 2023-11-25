package httpHeaderScan

import (
	"net/http"
	"strings"
)

func GenerateHeaderRecommendations(response *http.Response) []string {
	var (
		err                   error
		headerRecommendations = []string{}

		headerName  string
		headerValue string
	)

	// HSTS
	headerName = "Strict-Transport-Security"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		headerRecommendations = append(headerRecommendations, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#http-strict-transport-security")
	} else {
		err = validateHeaderStrictTransportSecurity(headerValue)
		if err != nil {
			headerRecommendations = append(headerRecommendations, "Recommended action for "+headerName+": "+err.Error())
		}
	}

	// Content-Security-Policy
	headerName = "Content-Security-Policy"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	if headerValue == "" {
		headerRecommendations = append(headerRecommendations, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#content-security-policy")
	} else {
		headerRecommendations = append(headerRecommendations, headerName+" header: "+headerValue) // TODO instead of just printing the header value check it against the best practices described in the link above
	}

	// X-Frame-Options
	headerName = "X-Frame-Options"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	switch headerValue {
	case "":
		headerRecommendations = append(headerRecommendations, headerName+" header should be set to 'sameorigin' or 'deny' as described in: https://infosec.mozilla.org/guidelines/web_security#x-frame-options")
	case "sameorigin":
		// Config okay
	case "deny":
		// Config okay
	default:
		headerRecommendations = append(headerRecommendations, headerName+" should be set to 'sameorigin' or 'deny', but got '"+headerValue+"'")
	}

	// X-Content-Type-Options
	headerName = "X-Content-Type-Options"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	switch headerValue {
	case "":
		headerRecommendations = append(headerRecommendations, headerName+" header should be set to 'nosniff' as described in: https://infosec.mozilla.org/guidelines/web_security#x-content-type-options")
	case "nosniff":
		// Perfectly configured
	default:
		headerRecommendations = append(headerRecommendations, headerName+" should be set to 'nosniff', but got '"+headerValue+"'")
	}

	// Referrer
	headerName = "Referer"
	headerValue = response.Header.Get(headerName)
	headerValue = strings.ToLower(headerValue)
	switch headerValue {
	case "":
		// headerRecommendations = append(headerRecommendations, headerName+" header should be implemented: https://infosec.mozilla.org/guidelines/web_security#referrer-policy")
		// Default is used, which is strict-origin-when-cross-origin and therefore okay
	case "no-referrer":
		// Config okay
	case "same-origin":
		// Config okay
	case "strict-origin":
		// Config okay
	case "strict-origin-when-cross-origin":
		// Config okay
	default:
		headerRecommendations = append(headerRecommendations, headerName+" should be set to one of 'no-referrer', 'same-origin', 'strict-origin', 'strict-origin-when-cross-origin', but got '"+headerValue+"'")
	}

	// Cache-Control
	headerName = "Cache-Control"
	headerValue = response.Header.Get(headerName)
	if headerValue == "" {
		headerRecommendations = append(headerRecommendations, headerName+" header should be configured, for example as described here: https://medium.com/pixelpoint/best-practices-for-cache-control-settings-for-your-website-ff262b38c5a2")
	} else {
		headerRecommendations = append(headerRecommendations, headerName+" header: "+headerValue) // TODO instead of just printing the header value check it against the best practices described in the link above
	}

	return headerRecommendations
}
