package htmlContentScan

import (
	"log/slog"

	"github.com/PuerkitoBio/goquery"
)

func GetExternalStylesheetSources(document *goquery.Document) []string {
	var (
		stylesheetSources = []string{}
	)

	slog.Debug("htmlContentScan: Getting external stylesheet sources started")

	// Find all `script` elements
	// <link rel="stylesheet" href="mystyle.css">
	document.Find("link").Each(func(i int, s *goquery.Selection) {
		if rel, exists := s.Attr("rel"); exists && rel == "stylesheet" {
			stylesheetSource, exists := s.Attr("href") // Get the `href` attribute for each `link` element with `rel="stylesheet` attribute set
			if exists {
				stylesheetSources = append(stylesheetSources, stylesheetSource)
			}
		}
	})

	slog.Debug("htmlContentScan: Getting external stylesheet sources completed")

	return stylesheetSources
}
