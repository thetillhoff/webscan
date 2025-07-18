package cachedHttpGetClient

import "time"

func (client Client) GetTimeout() time.Duration {
	return client.client.Timeout
}
