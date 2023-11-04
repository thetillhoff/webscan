package httpProtocolScan

import (
	"net/http"
	"net/url"
)

func checkHttp1(fullUrl string) (string, error) {
	var (
		err       error
		client    = &http.Client{}
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
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return response.Proto, nil
}
