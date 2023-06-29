package htmlcontent

import "github.com/PuerkitoBio/goquery"

func GetInlineStyleSize(document *goquery.Document) (int, error) {
	var (
		content string

		inlineStyleSize int = 0
	)

	// Find all `style` elements
	document.Find("style").Each(func(i int, s *goquery.Selection) {
		content, _ = s.Html() // Can't handle errors here
		inlineStyleSize = inlineStyleSize + len(content)
	})

	return inlineStyleSize, nil
}
