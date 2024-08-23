package htmlContentScan

type Result struct {
	httpContentHtmlSize         int
	httpContentInlineStyleSize  int
	httpContentInlineScriptSize int

	httpContentStylesheetSizes map[string]int
	httpContentScriptSizes     map[string]int

	httpContentRecommendations []string
}
