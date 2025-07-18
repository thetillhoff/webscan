package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/quic-go/quic-go/http3"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func checkHttp3(target types.Target) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http3.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // SSL verification is a different scan
				},
			}, // Required for http3
		}
		request  *http.Request
		response *http.Response
	)

	slog.Debug("httpProtocolScan: Checking http/3 started", "url", target.UrlString())

	if target.ParsedUrl().Scheme == "http" { // HTTP/3 is not supported for HTTP, so return fast
		err = errors.New("http/3 is not supported for HTTP")
		slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	// Create an HTTP request
	request, err = http.NewRequest("GET", target.UrlString(), nil)
	if err != nil {
		slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	// Add the Host header
	request.Header.Add("Host", target.Host()) // This is needed server-side to identify which vhost-config to use

	// Add the HTTP/3 header
	request.Header.Add("Alt-Svc", "h3=\":443\"")

	// Perform the HTTP request
	response, err = client.Do(request)

	if err == nil {
		defer func() {
			if closeErr := response.Body.Close(); closeErr != nil {
				slog.Debug("httpProtocolScan: Error closing response body", "error", closeErr)
			}
		}()

		if strings.HasPrefix(response.Proto, "HTTP/3") {
			slog.Debug("httpProtocolScan: Checking http/3 completed", "url", target.UrlString(), "proto", response.Proto)
			return response.Proto, nil
		} else {
			err = errors.New("http/3 is not supported")
			slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "proto", response.Proto, "error", err)
			return response.Proto, err
		}

	} else {
		slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "error", err)
		return "", err
	}
}
