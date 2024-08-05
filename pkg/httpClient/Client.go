package httpClient

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	client         *http.Client
	userAgent      string
	responses      map[string]*http.Response
	responseBodies map[string][]byte
}

func NewClient(timeout time.Duration, followRedirects int, verifyTls bool, userAgent string) Client {
	var client Client

	client.client = &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= followRedirects { // If amount of redirects is bigger or equal than set limit, don't follow further
				return http.ErrUseLastResponse
			} else { // If amount of redirects is lower than set limit, follow further
				return nil
			}
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTls},
		},
	}

	client.userAgent = userAgent
	client.responses = map[string]*http.Response{}
	client.responseBodies = map[string][]byte{}

	return client
}
