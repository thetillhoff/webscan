package webscan

import (
	"fmt"
	"strconv"
	"strings"
)

func (engine Engine) PrintHttpProtocolScanResults() {
	var (
		httpProtocolRecommendations = []string{}
	)

	if engine.HttpProtocolScan { // Only analyze protocol configuration if explicitly enabled

		if engine.isAvailableViaHttp { // If is available via http
			if engine.httpRedirectLocation != "" { // If http has redirectLocation
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP traffic is redirected to "+engine.httpRedirectLocation) // Print note with redirectLocation

				// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
				if engine.httpStatusCode != 301 && engine.httpStatusCode != 308 { // If status code != 301 or 308
					httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP should only be used to redirect with a 301 or 308 status code. Got "+strconv.Itoa(engine.httpStatusCode)) // Recommend to use 301 or 308 for redirect
				}

				if strings.HasPrefix(engine.httpRedirectLocation, "http://") { // If redirectLocation is Http
					httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP redirection targets should only be HTTPS locations. Got "+strconv.Itoa(engine.httpStatusCode)) // Recommend that http should only redirect to https location
				}
			} else if engine.httpStatusCode != 200 {
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP status code should be 200 when it's not used for redirects. Got "+strconv.Itoa(engine.httpStatusCode)) // Recommend that http should respond with 200 if it doesn't redirect
			}
		}

		if engine.isAvailableViaHttps { // If is available via https
			if engine.httpsRedirectLocation != "" { // If https has redirectLocation
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS traffic is redirected to "+engine.httpsRedirectLocation) // Print note with redirectLocation

				// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
				if engine.httpsStatusCode != 301 && engine.httpsStatusCode != 308 { // If status code != 301 or 308
					httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS should only be used to redirect with a 301 or 308 status code. Got "+strconv.Itoa(engine.httpsStatusCode)) // Recommend to use 301 or 308 for redirect
				}

				if strings.HasPrefix(engine.httpsRedirectLocation, "http://") { // If redirectLocation is Http
					httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS redirection targets should only be HTTPS locations. Got "+strconv.Itoa(engine.httpStatusCode)) // Recommend that https should only redirect to https location
				}
			} else if engine.httpsStatusCode != 200 {
				httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTPS status code should be 200 when it's not used for redirects. Got "+strconv.Itoa(engine.httpsStatusCode)) // Recommend that https should respond with 200 if it doesn't redirect
			}
		}

		if engine.isAvailableViaHttp && engine.isAvailableViaHttps { // If is available via Http and Https
			if engine.httpRedirectLocation != "" && engine.httpsRedirectLocation != "" { // If https is redirecting and http is redirecting
				if engine.httpRedirectLocation != engine.httpsRedirectLocation { // If http redirectLocation != https redirectLocation
					unifiedHttpRedirectLocation := strings.TrimSuffix(engine.httpRedirectLocation, "/")
					unifiedHttpsLocation := strings.TrimSuffix(engine.input, "/")
					unifiedHttpsLocation = strings.TrimPrefix(unifiedHttpsLocation, "https://")
					if unifiedHttpRedirectLocation != "https://"+unifiedHttpsLocation { // If http redirectLocation != self with https
						httpProtocolRecommendations = append(httpProtocolRecommendations, "HTTP and HTTPS are not redirecting to the same location, neither is HTTP redirecting to use HTTPS instead.") // Recommend to either redirect http to same target as https or just to https with same origin
					}
				}
			}
		}

		if len(httpProtocolRecommendations) > 0 {

			fmt.Printf("\n\n--- HTTP protocol scan results ---\n")

			for _, message := range httpProtocolRecommendations {
				fmt.Println(message)
			}
		}
	}

}
