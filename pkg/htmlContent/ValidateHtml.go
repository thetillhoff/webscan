package htmlcontent

import (
	"bytes"

	"golang.org/x/net/html"
)

// Only needs to be parsable, and to be html5

func ValidateHtml(body []byte) (string, error) {
	var (
		err      error
		doc      *html.Node
		htmlNode *html.Node
	)

	doc, err = html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	htmlNode = doc.FirstChild

	if htmlNode.Type != html.DoctypeNode { // First node is not doctype
		return "invalid html5: First non-comment node must be the doctype node, found " + htmlNode.Data, nil
	}

	for htmlNode.Type != html.ElementNode {
		htmlNode = htmlNode.NextSibling
		if htmlNode == nil {
			return "invalid html5: only doctype node found", nil
		}
	}

	if htmlNode.Data != "html" { // Next node must be the `html` node
		return "invalid html5: first documentnode after doctype must be `html` node, got " + htmlNode.Data, nil
	}

	// `html` node must have `lang` attr set
	langAttrExists := false
	for _, attr := range htmlNode.Attr {
		if attr.Key == "lang" {
			langAttrExists = true
			break
		}
	}
	if !langAttrExists {
		return "invalid html5: `html` node doesn't have `lang` attribute set", nil
	}

	return "", nil
}
