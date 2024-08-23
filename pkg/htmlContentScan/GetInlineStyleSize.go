package htmlContentScan

import (
	"log/slog"

	"github.com/PuerkitoBio/goquery"
)

func GetInlineStyleSize(document *goquery.Document) (int, error) {
	var (
		content string

		inlineStyleSize int = 0
	)

	slog.Debug("htmlContentScan: Getting inline style size started")

	// Find all `style` elements
	document.Find("style").Each(func(i int, s *goquery.Selection) {
		content, _ = s.Html() // Can't handle errors here
		inlineStyleSize = inlineStyleSize + len(content)
	})

	slog.Debug("htmlContentScan: Getting inline style size completed")

	return inlineStyleSize, nil
}
