package webscan

import (
	"fmt"
	"strconv"
	"strings"
)

func (engine Engine) PrintProtocolScanResults() {
	var (
		linebreakPrinted = false
	)

	if engine.HttpProtocolScan {

		if engine.isAvailableViaHttp && engine.httpRedirectLocation != "" { // If http does redirect
			if !linebreakPrinted {
				fmt.Println()
				linebreakPrinted = true
			}
			fmt.Println("HTTP traffic is redirected to", engine.httpRedirectLocation) // Display redirection location
		}
		if engine.isAvailableViaHttps && engine.httpsRedirectLocation != "" { // If https does redirect
			if !linebreakPrinted {
				fmt.Println()
				linebreakPrinted = true
			}
			fmt.Println("HTTPS traffic is redirected to", engine.httpsRedirectLocation) // Display redirect location
		}

		// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
		if engine.isAvailableViaHttp && engine.httpStatusCode != 301 && engine.httpStatusCode != 302 && engine.httpStatusCode != 303 && engine.httpStatusCode != 307 && engine.httpStatusCode != 308 { // Check against existing redirect status codes
			if !linebreakPrinted {
				fmt.Println()
				linebreakPrinted = true
			}
			fmt.Println("HTTP should only be used to redirect to an HTTPS location with a 301 or 308 status code. Got " + strconv.Itoa(engine.httpStatusCode))
		}

		if engine.isAvailableViaHttps && engine.httpsStatusCode != 200 {
			if !linebreakPrinted {
				fmt.Println()
				linebreakPrinted = true
			}
			fmt.Println("HTTPS status code should be 200 when it's not used for redirects. Got " + strconv.Itoa(engine.httpsStatusCode))
		}

		if engine.httpsRedirectLocation != "" { // If https redirects
			if engine.httpRedirectLocation != engine.httpsRedirectLocation { // Http and https should redirect to exact same location
				if strings.TrimSuffix(engine.httpRedirectLocation, "/") != "https://"+engine.url { // http does not redirect to https (same origin)
					if !linebreakPrinted {
						fmt.Println()
						linebreakPrinted = true
					}
					fmt.Println("Both HTTP and HTTPS are redirecting, so they should redirect to the same location. Instead got " + engine.httpRedirectLocation + " for http and " + engine.httpsRedirectLocation + " for https")
				}
			}

			if !strings.HasPrefix(engine.httpRedirectLocation, "https://") { // Not redirecting to  https location
				if !linebreakPrinted {
					fmt.Println()
					linebreakPrinted = true
				}
				fmt.Println("Both HTTP and HTTPS are redirecting, and should redirect to https locations only. Instead got: " + engine.httpsRedirectLocation)
			}
		} else { // Else https serves a page

			if engine.isAvailableViaHttp && engine.httpStatusCode != 301 && engine.httpStatusCode != 308 { // Http should redirect to https with 301 or 308
				if !linebreakPrinted {
					fmt.Println()
					linebreakPrinted = true
				}
				fmt.Println("HTTP redirect to HTTPS should happen with 301 or 308 status code. Instead got: " + strconv.Itoa(engine.httpStatusCode))
			}

			if engine.isAvailableViaHttp && strings.TrimSuffix(engine.httpRedirectLocation, "/") != "https://"+engine.url { // http does not redirect to https (same origin)
				if !linebreakPrinted {
					fmt.Println()
					linebreakPrinted = true
				}
				fmt.Println("HTTP should redirect to HTTPS. Instead got: " + engine.httpRedirectLocation)
			}
		}

		// Both cases for isAvailableViaHttp are covered now, so we know the page is either accessible via https or not at all.
		// -> All is good.
	}
}
