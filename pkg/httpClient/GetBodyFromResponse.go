package httpClient

import (
	"io"
	"net/http"
)

func GetBodyFromResponse(response *http.Response) ([]byte, error) {
	var (
		err error

		body []byte
	)

	body, err = io.ReadAll(response.Body) // Read the response
	if err != nil {
		return body, err
	}
	defer response.Body.Close()

	return body, nil
}
