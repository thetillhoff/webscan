package protocolScan

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/quic-go/quic-go/http3"
)

func checkHttp3(fullUrl string) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http3.RoundTripper{}, // Required for http2
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

	// Perform the HTTP request
	response, err = client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check if the response indicates HTTP/3 support
	if response.ProtoMajor == 3 {
		fmt.Println("HTTP/3 (QUIC) is supported")
		return response.Proto, nil
	} else {
		fmt.Println("HTTP/3 (QUIC) is not supported")
		return "", nil
	}
}
