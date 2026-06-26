package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/types"
)

func checkHTTP1(target types.Target, timeout time.Duration) (string, error) {
	var (
		err    error
		client = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				ForceAttemptHTTP2: false,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0),
			},
		}
		request  *http.Request
		response *http.Response
	)

	slog.Debug("httpProtocolScan: Checking http/1 started", "url", target.UrlString())

	request, err = http.NewRequest("GET", target.UrlString(), nil)
	if err != nil {
		slog.Debug("httpProtocolScan: Checking http/1 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	request.Header.Add("Host", target.Host())

	response, err = client.Do(request)
	if err == nil {
		defer func() {
			if closeErr := response.Body.Close(); closeErr != nil {
				slog.Debug("httpProtocolScan: Error closing response body", "error", closeErr)
			}
		}()

		if strings.HasPrefix(response.Proto, "HTTP/1") {
			slog.Debug("httpProtocolScan: Checking http/1 completed", "url", target.UrlString(), "proto", response.Proto)
			return response.Proto, nil
		} else {
			err = errors.New("http/1 is not supported")
			slog.Debug("httpProtocolScan: Checking http/1 failed", "url", target.UrlString(), "proto", response.Proto, "error", err)
			return response.Proto, err
		}

	} else {
		slog.Debug("httpProtocolScan: Checking http/1 failed", "url", target.UrlString(), "error", err)
		return "", err
	}
}
