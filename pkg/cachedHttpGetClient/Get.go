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

	if cachedResponse, ok := client.responses[url]; ok { // If cached response exists

		slog.Debug("httpClient: Returning response for request from internal cache", "url", url)

		return cachedResponse.GetHttpResponse(), cachedResponse.GetBody(), cachedResponse.GetError()

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

		client.responses[url] = response

		slog.Debug("httpClient: Request completed", "url", url)

		return response.GetHttpResponse(), response.GetBody(), response.GetError()
	}
}
