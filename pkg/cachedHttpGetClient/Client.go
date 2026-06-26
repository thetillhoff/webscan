package cachedHttpGetClient

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	client    *http.Client
	userAgent string
	verifyTls bool
	mu        *sync.Mutex
	// ponytail: cache is unbounded but scoped to one short-lived engine instance per scan; upgrade to an LRU if engines become long-lived.
	responses map[string]Response
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
	client.verifyTls = verifyTls
	client.mu = &sync.Mutex{}
	client.responses = map[string]Response{}

	return client
}
