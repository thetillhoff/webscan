package httpClient

import (
	"io"
	"log/slog"
	"net/http"
)

func (client Client) MakeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	var (
		err error

		request  *http.Request
		response *http.Response
	)

	slog.Debug("httpClient: HttpRequest requested", "method", method, "url", url)

	if cachedResponse, ok := client.responses[method+url]; ok { // If cached response exists

		slog.Debug("httpClient: Returning response for request from internal cache", "method", method, "url", url)

		return cachedResponse, nil

	} else { // If no cached response exists

		slog.Debug("httpClient: Making request", "method", method, "url", url)

		request, err = http.NewRequest(method, url, body)
		if err != nil {
			return response, err
		}
		request.Header.Set("User-Agent", client.userAgent) // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
		response, err = client.client.Do(request)
		if err != nil {
			return response, err
		}

		client.responses[method+url] = response // Add response to cache

		slog.Debug("httpClient: Request completed")

		return response, nil
	}
}
