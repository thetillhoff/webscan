package httpProtocolScan

type Result struct {
	isAvailableViaHttp   bool
	httpStatusCode       int
	httpRedirectLocation string
	httpVersions         []string

	isAvailableViaHttps   bool
	httpsStatusCode       int
	httpsRedirectLocation string
	httpsVersions         []string

	recommendations []string
}

func (result Result) IsAvailableViaHttp() bool {
	return result.isAvailableViaHttp
}

func (result Result) IsAvailableViaHttps() bool {
	return result.isAvailableViaHttps
}
