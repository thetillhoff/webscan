package cachedHttpGetClient

func (client Client) VerifyTLS() bool {
	return client.verifyTls
}
