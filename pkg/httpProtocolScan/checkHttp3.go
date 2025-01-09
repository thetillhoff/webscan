package httpProtocolScan

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/quic-go/quic-go/http3"
)

func checkHttp3(fullUrl string) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http3.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // SSL verification is a different scan
				},
			}, // Required for http3
		}
		parsedUrl *url.URL
		request   *http.Request
		response  *http.Response
	)

	parsedUrl, err = url.Parse(fullUrl)
	if err != nil {
		return "", err
	}

	// Create an HTTP request
	request, err = http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return "", err
	}

	// Add the Host header
	request.Header.Add("Host", parsedUrl.Host) // This is needed server-side to identify which vhost-config to use

	// Add the HTTP/3 header
	request.Header.Add("Alt-Svc", "h3=\":443\"")

	// Perform the HTTP request
	response, err = client.Do(request)

	if err == nil {
		defer response.Body.Close()
		slog.Debug("Result of check for http/3 protocol support", "proto", response.Proto)
		return response.Proto, nil
	} else {
		return "", nil
	}
}
