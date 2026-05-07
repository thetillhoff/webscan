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

func (result Result) HttpRedirectLocation() string {
	return result.httpRedirectLocation
}

func (result Result) HttpsRedirectLocation() string {
	return result.httpsRedirectLocation
}

func (result Result) HttpStatusCode() int {
	return result.httpStatusCode
}

func (result Result) HttpsStatusCode() int {
	return result.httpsStatusCode
}
