package httpProtocolScan

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func checkHTTP3(target types.Target, timeout time.Duration) (string, error) {
	var (
		err    error
		client = &http.Client{
			Timeout: timeout,
			Transport: &http3.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		request  *http.Request
		response *http.Response
	)

	slog.Debug("httpProtocolScan: Checking http/3 started", "url", target.UrlString())

	if target.ParsedUrl().Scheme == "http" {
		err = errors.New("http/3 is not supported for HTTP")
		slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	request, err = http.NewRequest("GET", target.UrlString(), nil)
	if err != nil {
		slog.Debug("httpProtocolScan: Checking http/3 failed", "url", target.UrlString(), "error", err)
		return "", err
	}

	request.Header.Add("Host", target.Host())
	request.Header.Add("Alt-Svc", "h3=\":443\"")

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
