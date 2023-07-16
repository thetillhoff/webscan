package tlsScan

import (
	"errors"
	"net/http"
	"os"
	"time"
)

// checks whether the certificate is valid
func ValidateTlsCertificate(url string) error {
	var (
		err error
	)

	client := &http.Client{
		Timeout: 5 * time.Second, // TODO 5s might be a bit long?
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // Don't follow redirects
	}

	_, err = client.Get("https://" + url)

	if err != nil {
		if os.IsTimeout(err) { // If err is timeout, tell user about it
			return errors.New("http call exceeded 5s timeout")
		} else {
			return errors.New("Invalid TLS certificate used for HTTPS: " + errors.Unwrap(errors.Unwrap(err)).Error())
		}
	}

	return nil
}
