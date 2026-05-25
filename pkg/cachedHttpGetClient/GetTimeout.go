package cachedHttpGetClient

import "time"

func (client Client) Timeout() time.Duration {
	return client.client.Timeout
}
