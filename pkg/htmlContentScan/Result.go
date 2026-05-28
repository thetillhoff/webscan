package htmlContentScan

import "reflect"

type Result struct {
	httpContentHtmlSize         int
	httpContentInlineStyleSize  int
	httpContentInlineScriptSize int

	httpContentStylesheetSizes map[string]int
	httpContentScriptSizes     map[string]int

	httpContentRecommendations []string
}

func (r Result) Equal(other Result) bool {
	return r.httpContentHtmlSize == other.httpContentHtmlSize &&
		r.httpContentInlineStyleSize == other.httpContentInlineStyleSize &&
		r.httpContentInlineScriptSize == other.httpContentInlineScriptSize &&
		reflect.DeepEqual(r.httpContentStylesheetSizes, other.httpContentStylesheetSizes) &&
		reflect.DeepEqual(r.httpContentScriptSizes, other.httpContentScriptSizes) &&
		reflect.DeepEqual(r.httpContentRecommendations, other.httpContentRecommendations)
}
