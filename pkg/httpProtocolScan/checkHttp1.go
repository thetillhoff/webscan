package httpProtocolScan

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/url"
)

func checkHttp1(fullUrl string) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http.Transport{
				ForceAttemptHTTP2: false, // Disable HTTP/2
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // SSL verification is a different scan
				},
				TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0), // Disable HTTP/2
			},
		}
		parsedUrl *url.URL
		request   *http.Request
		response  *http.Response
	)

	parsedUrl, err = url.Parse(fullUrl)
	if err != nil {
		return "", err
	}

	request, err = http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Host", parsedUrl.Host) // This is needed server-side to identify which vhost-config to use

	response, err = client.Do(request)
	if err == nil {
		defer response.Body.Close()
		slog.Debug("Result of check for http/1.1 protocol support", "proto", response.Proto)
		return response.Proto, nil
	} else {
		return "", nil
	}
}
