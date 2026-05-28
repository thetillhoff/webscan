package knownFilesScan

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

func analyzeSitemapXml(body []byte) (observations []string, recommendations []string) {
	var (
		urlCount       int
		sitemapCount   int
		isSitemapIndex bool
	)

	decoder := xml.NewDecoder(bytes.NewReader(body))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			recommendations = append(recommendations, "Could not parse sitemap.xml as valid XML: "+err.Error())
			return observations, recommendations
		}

		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "sitemapindex":
			isSitemapIndex = true
		case "url":
			urlCount++
		case "sitemap":
			if isSitemapIndex {
				sitemapCount++
			}
		}
	}

	if isSitemapIndex {
		observations = append(observations, fmt.Sprintf("Sitemap index with %d child sitemaps", sitemapCount))
	} else if urlCount > 0 {
		observations = append(observations, fmt.Sprintf("%d URLs listed", urlCount))
		if urlCount > 50000 {
			recommendations = append(recommendations, fmt.Sprintf("Exceeds the 50,000 URL limit (%d URLs); split into a sitemap index", urlCount))
		}
	} else {
		recommendations = append(recommendations, "No <url> or <sitemap> elements found")
	}

	return observations, recommendations
}
