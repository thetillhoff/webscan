package httpProtocolScan

type Result struct {
	httpStatusCode       int
	httpRedirectLocation string
	httpVersions         []string

	httpsStatusCode       int
	httpsRedirectLocation string
	httpsVersions         []string

	recommendations []string
}
