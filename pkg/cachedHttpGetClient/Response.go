package cachedHttpGetClient

import (
	"io"
	"net/http"
)

// ponytail: cap response bodies at 10 MB to prevent memory exhaustion from hostile targets; raise if larger bodies must be scanned.
const maxBodySize = 10 << 20 // 10 MB

type Response struct {
	httpResponse *http.Response
	body         []byte
	err          error
}

func NewResponse(httpResponse *http.Response, err error) Response {

	if err != nil {
		return Response{
			httpResponse: httpResponse,
			body:         nil,
			err:          err,
		}
	}

	body, err := io.ReadAll(io.LimitReader(httpResponse.Body, maxBodySize))
	if closeErr := httpResponse.Body.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	return Response{
		httpResponse: httpResponse,
		body:         body,
		err:          err,
	}
}

func (response Response) HTTPResponse() *http.Response {
	return response.httpResponse
}

func (response Response) Body() []byte {
	return response.body
}

func (response Response) Err() error {
	return response.err
}
