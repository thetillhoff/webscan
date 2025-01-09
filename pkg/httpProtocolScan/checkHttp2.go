package httpProtocolScan

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

func checkHttp2(fullUrl string) (string, error) {
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
		fmt.Println("proto2", response.Proto)
		return response.Proto, nil
	} else {
		return "", nil
	}
}
