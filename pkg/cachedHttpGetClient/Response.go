package cachedHttpGetClient

import (
	"io"
	"net/http"
)

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

	body, err := io.ReadAll(httpResponse.Body)
	return Response{
		httpResponse: httpResponse,
		body:         body,
		err:          err,
	}
}

func (response Response) GetHttpResponse() *http.Response {
	return response.httpResponse
}

func (response Response) GetBody() []byte {
	return response.body
}

func (response Response) GetError() error {
	return response.err
}
