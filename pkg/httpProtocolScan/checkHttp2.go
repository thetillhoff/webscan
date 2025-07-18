package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func checkHttp2(target types.Target) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http.Transport{
				ForceAttemptHTTP2: true, // Force HTTP/2
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // SSL verification is a different scan
				},
			}, // Required for http2
		}
		request  *http.Request
		response *http.Response
	)

	slog.Debug("httpProtocolScan: Checking http/2 started", "url", target.UrlString())

	request, err = http.NewRequest("GET", target.UrlString(), nil)
	if err != nil {
		slog.Debug("httpProtocolScan: Checking http/2 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	request.Header.Add("Host", target.Host()) // This is needed server-side to identify which vhost-config to use

	response, err = client.Do(request)
	if err == nil {
		defer func() {
			if closeErr := response.Body.Close(); closeErr != nil {
				slog.Debug("httpProtocolScan: Error closing response body", "error", closeErr)
			}
		}()

		if strings.HasPrefix(response.Proto, "HTTP/2") {
			slog.Debug("httpProtocolScan: Checking http/2 completed", "url", target.UrlString(), "proto", response.Proto)
			return response.Proto, nil
		} else {
			err = errors.New("http/2 is not supported")
			slog.Debug("httpProtocolScan: Checking http/1 failed", "url", target.UrlString(), "proto", response.Proto, "error", err)
			return response.Proto, err
		}

	} else {
		slog.Debug("httpProtocolScan: Checking http/2 failed", "url", target.UrlString(), "error", err)
		return "", err
	}
}
