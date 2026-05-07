package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// CheckHTTPRedirects checks if the target redirects without following it.
// Returns the status code, redirect location (empty if no redirect), and error.
func CheckHTTPRedirects(target types.Target, timeout time.Duration) (int, string, error) {
	var (
		statusCode       = 0
		redirectLocation = ""
	)

	slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects started", target.Schema().String()))

	noFollowClient := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := noFollowClient.Get(target.UrlString())
	if err != nil {
		if os.IsTimeout(err) {
			return statusCode, redirectLocation, errors.New("http call exceeded timeout")
		}
		slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects unsuccessful", target.Schema().String()), "error", err)
		return statusCode, redirectLocation, err
	}

	defer func() { _ = resp.Body.Close() }()

	statusCode = resp.StatusCode
	if statusCode == 301 || statusCode == 302 || statusCode == 303 || statusCode == 307 || statusCode == 308 {
		redirectLocation = resp.Header.Get("Location")
	}

	slog.Debug(fmt.Sprintf("httpProtocolScan: Checking %s redirects completed", target.Schema().String()), "statusCode", statusCode, "redirectLocation", redirectLocation)

	return statusCode, redirectLocation, nil
}
