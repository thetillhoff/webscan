package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"time"
)

// Verifies http and https configuration for best practices
// Returns http and https status codes & redirect locations, which are "" if they don't redirect plus an potential error
func CheckHttpRedirects(url string, isAvailableViaHttp bool, isAvailableViaHttps bool) (int, string, int, string, error) {
	var (
		err       error
		httpResp  *http.Response
		httpsResp *http.Response

		httpStatusCode       int    = 0
		httpRedirectLocation string = "" // Might contain a trailing '/'

		httpsStatusCode       int    = 0
		httpsRedirectLocation string = "" // Might contain a trailing '/'
	)

	client := &http.Client{
		Timeout: 5 * time.Second, // TODO 5s might be a bit long?
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // Don't follow redirects
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore invalid tls certificates here (certificates are checked in another step, and might be interesting what's behind it anyway)
		},
	}

	if isAvailableViaHttp { // If HTTP is available, check the status code
		// Get the page in HTTP
		httpResp, err = client.Get("http://" + url)
		if err != nil {
			if os.IsTimeout(err) { // If err is timeout, tell user about it
				return httpStatusCode, httpRedirectLocation, httpsStatusCode, httpsRedirectLocation, errors.New("http call exceeded 5s timeout")
			} else {
				return httpStatusCode, httpRedirectLocation, httpsStatusCode, httpsRedirectLocation, err
			}
		}

		defer httpResp.Body.Close()

		// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
		httpStatusCode = httpResp.StatusCode                                                                                           // Get status code
		if httpStatusCode == 301 || httpStatusCode == 302 || httpStatusCode == 303 || httpStatusCode == 307 || httpStatusCode == 308 { // Check against existing redirect status codes
			httpRedirectLocation = httpResp.Header.Get("Location") // Get Location
		}
	}

	if isAvailableViaHttps { // If HTTPS is available, check the status code
		// Get the page in HTTPS
		httpsResp, err = client.Get("https://" + url)
		if err != nil {
			if os.IsTimeout(err) { // If err is timeout, tell user about it
				return httpStatusCode, httpRedirectLocation, httpsStatusCode, httpsRedirectLocation, errors.New("https call exceeded 5s timeout")
			} else {
				return httpStatusCode, httpRedirectLocation, httpsStatusCode, httpsRedirectLocation, err
			}
		}

		defer httpsResp.Body.Close()

		// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases not meant for normal redirects
		httpsStatusCode = httpsResp.StatusCode                                                                                          // Get status code
		if httpsStatusCode == 301 || httpStatusCode == 302 || httpStatusCode == 303 || httpStatusCode == 307 || httpStatusCode == 308 { // Check against existing redirect status codes
			httpsRedirectLocation = httpsResp.Header.Get("Location") // Get Location
		}
	}

	return httpStatusCode, httpRedirectLocation, httpsStatusCode, httpsRedirectLocation, nil
}
