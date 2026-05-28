package httpHeaderScan

type HeaderEntry struct {
	Name           string
	Value          string // empty means not set; may contain \n for multi-line values (e.g. CSP)
	Recommendation string // empty means no recommendation
}

type Result struct {
	httpHeaderEntries              []HeaderEntry
	httpCookieRecommendations      map[string][]string
	httpOtherCookieRecommendations []string
}
