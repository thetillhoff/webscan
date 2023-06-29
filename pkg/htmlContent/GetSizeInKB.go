package htmlcontent

import (
	"io"
	"math"
	"net/http"
	"time"
)

func GetSizeInKb(url string) (float32, error) {
	var (
		err error

		client   *http.Client
		request  *http.Request
		response *http.Response

		body []byte
		size float64
	)

	client = &http.Client{
		Timeout: 5 * time.Second, // TODO 5s might be a bit long?
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // Don't follow redirects // TODO Should we follow redirects or not?
	}
	request, err = http.NewRequest("GET", url, nil) // Only for https pages.
	if err != nil {
		return 0, err
	}
	request.Header.Set("User-Agent", "Go-http-client/1.1") // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
	// TODO request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0") // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
	response, err = client.Do(request)
	if err != nil {
		return 0, err
	}

	body, err = io.ReadAll(response.Body) // Read the response
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	size = float64(len(body))
	size = size / 1000              // Byte to Kilobyte conversion
	size = math.Round(size*10) / 10 // Rounding like this will leave exactly one decimal

	return float32(size), nil
}
