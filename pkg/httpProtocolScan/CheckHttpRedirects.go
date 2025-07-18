package httpProtocolScan

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/thetillhoff/webscan/v3/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// Verifies http and https configuration for best practices
// Returns http and https status codes & redirect locations, which are "" if they don't redirect plus an potential error
func CheckHttpRedirects(target types.Target, client cachedHttpGetClient.Client) (int, string, error) {
	var (
		err  error
		resp *http.Response

		statusCode       = 0
		redirectLocation = "" // Might contain a trailing '/'
	)

	slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects started", target.Schema().String()))

	resp, _, err = client.Get(target.UrlString())

	if err != nil {
		if os.IsTimeout(err) { // If err is timeout, tell user about it
			return statusCode, redirectLocation, errors.New("http call exceeded 5s timeout")
		} else {
			slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects unsuccessful", target.Schema().String()), "error", err)
			return statusCode, redirectLocation, err
		}
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Debug("httpProtocolScan: Error closing response body", "error", closeErr)
		}
	}()

	// 301 & 308 are permanent redirects, 302, 303, 307 are temporary redirects, 300 and 304 are special cases are not meant for normal redirects
	statusCode = resp.StatusCode                                                                               // Get status code
	if statusCode == 301 || statusCode == 302 || statusCode == 303 || statusCode == 307 || statusCode == 308 { // Check against existing redirect status codes
		redirectLocation = resp.Header.Get("Location") // Get Location
	}

	slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects completed", target.Schema().String()), "statusCode", statusCode, "redirectLocation", redirectLocation)

	return statusCode, redirectLocation, nil
}
