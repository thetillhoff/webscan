package httpProtocolScan

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func GetHttpProtocolRecommendationsForResult(input string, result Result, isAvailableViaHttp bool, isAvailableViaHttps bool) []string {
	var (
		httpProtocolRecommendations = []string{}
	)

	slog.Debug("httpProtocolScan: Getting recommendations for http protocol result started")

	if isAvailableViaHttp { // If is available via http
		if result.httpRedirectLocation != "" { // If http has redirectLocation
			httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP traffic is redirected to "+result.httpRedirectLocation) // Print note with redirectLocation

			// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
			if result.httpStatusCode != 301 && result.httpStatusCode != 308 { // If status code != 301 or 308
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP should only be used to redirect with a 301 or 308 status code. Got "+strconv.Itoa(result.httpStatusCode)) // Recommend to use 301 or 308 for redirect
			}

			if strings.HasPrefix(result.httpRedirectLocation, "http://") { // If redirectLocation is Http
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP redirection targets should only be HTTPS locations. Got "+strconv.Itoa(result.httpStatusCode)) // Recommend that http should only redirect to https location
			}
		} else if result.httpStatusCode != 200 {
			httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP status code should be 200 when it's not used for redirects. Got "+strconv.Itoa(result.httpStatusCode)) // Recommend that http should respond with 200 if it doesn't redirect
		}
	}

	if isAvailableViaHttps { // If is available via https
		if result.httpsRedirectLocation != "" { // If https has redirectLocation
			httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS traffic is redirected to "+result.httpsRedirectLocation) // Print note with redirectLocation

			// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
			if result.httpsStatusCode != 301 && result.httpsStatusCode != 308 { // If status code != 301 or 308
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS should only be used to redirect with a 301 or 308 status code. Got "+strconv.Itoa(result.httpsStatusCode)) // Recommend to use 301 or 308 for redirect
			}

			if strings.HasPrefix(result.httpsRedirectLocation, "http://") { // If redirectLocation is Http
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS redirection targets should only be HTTPS locations. Got "+strconv.Itoa(result.httpStatusCode)) // Recommend that https should only redirect to https location
			}
		} else if result.httpsStatusCode != 200 {
			httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS status code should be 200 when it's not used for redirects. Got "+strconv.Itoa(result.httpsStatusCode)) // Recommend that https should respond with 200 if it doesn't redirect
		}
	}

	if isAvailableViaHttp && isAvailableViaHttps { // If is available via Http and Https
		if result.httpRedirectLocation != "" && result.httpsRedirectLocation != "" { // If https is redirecting and http is redirecting
			slog.Info("HTTP and HTTPS are both redirecting somewhere")
			if result.httpRedirectLocation != result.httpsRedirectLocation { // If http redirectLocation != https redirectLocation
				slog.Info("HTTP and HTTPS are not redirecting to the same location")

				if !strings.HasPrefix(result.httpRedirectLocation, fmt.Sprintf("https://%s", input)) { // If httpRedirectLocation starts with 'https://<target>'
					httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP and HTTPS are not redirecting to the same location, neither is HTTP redirecting to use HTTPS instead.") // Recommend to either redirect http to same target as https or just to https with same origin
				}
			}
		}
	}

	slog.Debug("httpProtocolScan: Getting recommendations for http protocol result completed")

	return httpProtocolRecommendations
}
