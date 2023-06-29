package protocolScan

import (
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

func checkHttp2(fullUrl string) (string, error) {
	var (
		err    error
		client = &http.Client{
			Transport: &http2.Transport{}, // Required for http2
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
		return response.Proto, nil
	} else {
		return "", nil
	}
}
