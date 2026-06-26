package cachedHttpGetClient

import (
	"log/slog"
	"net/http"
)

func (client Client) Get(url string) (*http.Response, []byte, error) {
	var (
		err error

		request      *http.Request
		httpResponse *http.Response
		response     Response
	)

	slog.Debug("httpClient: HttpRequest requested", "url", url)

	client.mu.Lock()
	cachedResponse, ok := client.responses[url]
	client.mu.Unlock()

	if ok { // If cached response exists

		slog.Debug("httpClient: Returning response for request from internal cache", "url", url)

		return cachedResponse.HTTPResponse(), cachedResponse.Body(), cachedResponse.Err()

	} else { // If no cached response exists

		slog.Debug("httpClient: Making request", "url", url)

		request, err = http.NewRequest("GET", url, nil)
		if err != nil {
			response = NewResponse(nil, err)
		} else {
			request.Header.Set("User-Agent", client.userAgent) // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
			httpResponse, err = client.client.Do(request)

			response = NewResponse(httpResponse, err)
		}

		client.mu.Lock()
		client.responses[url] = response
		client.mu.Unlock()

		slog.Debug("httpClient: Request completed", "url", url)

		return response.HTTPResponse(), response.Body(), response.Err()
	}
}
