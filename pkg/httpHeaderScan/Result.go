package httpHeaderScan

type Result struct {
	httpHeaderRecommendations      []string
	httpCookieRecommendations      map[string][]string
	httpOtherCookieRecommendations []string
}
