package knownFilesScan

import (
	"reflect"
	"sort"

	"github.com/thetillhoff/webscan/v5/pkg/types"
)

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

// EqualContent returns true when both results have the same file findings (ignoring schema).
func (r Result) EqualContent(other Result) bool {
	if len(r.files) != len(other.files) {
		return false
	}
	rFiles := make([]FileResult, len(r.files))
	oFiles := make([]FileResult, len(other.files))
	copy(rFiles, r.files)
	copy(oFiles, other.files)
	sort.Slice(rFiles, func(i, j int) bool { return rFiles[i].Path < rFiles[j].Path })
	sort.Slice(oFiles, func(i, j int) bool { return oFiles[i].Path < oFiles[j].Path })
	for i := range rFiles {
		if rFiles[i].Path != oFiles[i].Path ||
			rFiles[i].Found != oFiles[i].Found ||
			rFiles[i].StatusCode != oFiles[i].StatusCode ||
			!reflect.DeepEqual(rFiles[i].Observations, oFiles[i].Observations) ||
			!reflect.DeepEqual(rFiles[i].Recommendations, oFiles[i].Recommendations) {
			return false
		}
	}
	return true
}
