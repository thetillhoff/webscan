package httpClient

import (
	"io"
	"net/http"
)

func (httpClient Client) MakeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	var (
		err error

		request  *http.Request
		response *http.Response
	)

	request, err = http.NewRequest(method, url, body) // Only for https pages.
	if err != nil {
		return response, err
	}
	request.Header.Set("User-Agent", httpClient.userAgent) // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
	response, err = httpClient.client.Do(request)
	if err != nil {
		return response, err
	}

	return response, nil
}
