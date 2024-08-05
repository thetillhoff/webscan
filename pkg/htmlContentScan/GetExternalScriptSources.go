package htmlContentScan

import (
	"log/slog"

	"github.com/PuerkitoBio/goquery"
)

func GetExternalScriptSources(document *goquery.Document) []string {
	var (
		scriptSources = []string{}
	)

	slog.Debug("htmlContentScan: Getting external script sources started")

	// Find all `script` elements
	document.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptSrc, exists := s.Attr("src") // Get the `src` attribute for each `script` element
		if exists {
			scriptSources = append(scriptSources, scriptSrc)
		}
	})

	slog.Debug("htmlContentScan: Getting external script sources completed")

	return scriptSources
}
