package htmlContentScan

import (
	"log/slog"

	"github.com/PuerkitoBio/goquery"
)

func GetInlineScriptSize(document *goquery.Document) (int, error) {
	var (
		content string

		inlineScriptSize = 0
	)

	slog.Debug("htmlContentScan: Getting inline script size started")

	// Find all `script` elements
	document.Find("script").Each(func(i int, s *goquery.Selection) {
		_, exists := s.Attr("src") // Get the `src` attribute for each `script` element
		if !exists {
			content, _ = s.Html() // Can't handle errors here
			inlineScriptSize = inlineScriptSize + len(content)
		}
	})

	slog.Debug("htmlContentScan: Getting inline script size completed")

	return inlineScriptSize, nil
}
