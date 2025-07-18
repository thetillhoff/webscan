package cachedHttpGetClient

func (client Client) DoesVerifyTls() bool {
	return client.verifyTls
}
