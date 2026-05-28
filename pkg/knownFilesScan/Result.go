package knownFilesScan

import "github.com/thetillhoff/webscan/v3/pkg/types"

type FileResult struct {
	Path            string
	Label           string
	Category        fileCategory
	StatusCode      int
	Found           bool
	Observations    []string
	Recommendations []string
	sitemapURLs     []string // extracted from robots.txt for cross-referencing
}

type Result struct {
	schema types.Schema
	files  []FileResult
}
