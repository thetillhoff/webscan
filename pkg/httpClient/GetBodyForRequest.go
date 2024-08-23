package httpClient

import (
	"io"
	"net/http"
)

func (client Client) GetBodyForRequest(method string, url string) ([]byte, error) {
	var (
		err error

		response *http.Response
		body     []byte
	)

	if cachedResponseBody, ok := client.responseBodies[method+url]; ok { // If cached response exists

		return cachedResponseBody, nil

	} else {

		response, err = client.MakeRequest(method, url, nil) // Make request
		if err != nil {
			return body, err
		}

		body, err = io.ReadAll(response.Body) // Read body from response
		if err != nil {
			return body, err
		}
		defer response.Body.Close()

		client.responseBodies[method+url] = body // Add body to cache

		return body, nil
	}

}
