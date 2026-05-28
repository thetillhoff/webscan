package knownFilesScan

import (
	"bufio"
	"bytes"
	"strings"
)

func analyzeRobotsTxt(body []byte) (observations []string, recommendations []string, sitemapURLs []string) {
	var (
		inStarBlock bool
	)

	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "user-agent":
			inStarBlock = value == "*"
		case "disallow":
			if inStarBlock && value == "/" {
				recommendations = append(recommendations, "'Disallow: /' under 'User-agent: *' blocks all crawlers from the entire site")
			}
		case "sitemap":
			sitemapURLs = append(sitemapURLs, value)
		}
	}

	if len(sitemapURLs) > 0 {
		for _, u := range sitemapURLs {
			observations = append(observations, "Sitemap directive: "+u)
		}
	} else {
		recommendations = append(recommendations, "No Sitemap directive found; add 'Sitemap: <url>' to help crawlers discover your sitemap")
	}

	return observations, recommendations, sitemapURLs
}
