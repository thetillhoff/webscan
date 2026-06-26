package knownFilesScan

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

type fileCategory int

const (
	categoryExpected  fileCategory = iota
	categorySensitive fileCategory = iota
)

type fileCheck struct {
	path     string
	label    string
	category fileCategory
}

var knownFiles = []fileCheck{
	{"/robots.txt", "robots.txt", categoryExpected},
	{"/sitemap.xml", "sitemap.xml", categoryExpected},
	{"/.well-known/security.txt", "security.txt", categoryExpected},
	{"/llms.txt", "llms.txt", categoryExpected},
	{"/.well-known/ai-plugin.json", "AI plugin manifest", categoryExpected},

	{"/.htaccess", ".htaccess", categorySensitive},
	{"/.env", ".env", categorySensitive},
	{"/.git/config", ".git/config", categorySensitive},
	{"/wp-config.php", "wp-config.php", categorySensitive},
	{"/server-status", "server-status", categorySensitive},
}

type scanConfig struct {
	timeout time.Duration
}

type ConfigOption func(*scanConfig)

func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

func Scan(target types.Target, status *status.Status, schema types.Schema, options ...ConfigOption) Result {
	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	for _, option := range options {
		option(config)
	}

	slog.Debug("knownFilesScan: Scan started")

	status.SpinningUpdate(fmt.Sprintf("Scanning %s well-known files...", schema.String()))

	result := Result{
		schema: schema,
	}

	httpClient := &http.Client{
		Timeout: config.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, fc := range knownFiles {
		wg.Add(1)
		go func(fc fileCheck) {
			defer wg.Done()

			url := schema.String() + "://" + target.Hostname() + fc.path

			resp, err := httpClient.Get(url)
			if err != nil {
				slog.Debug("knownFilesScan: request failed", "url", url, "error", err)
				return
			}
			defer func() { _ = resp.Body.Close() }()

			entry := FileResult{
				Path:       fc.path,
				Label:      fc.label,
				Category:   fc.category,
				StatusCode: resp.StatusCode,
				Found:      resp.StatusCode == http.StatusOK,
			}

			if entry.Found {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					slog.Debug("knownFilesScan: failed to read body", "url", url, "error", err)
				} else {
					switch fc.path {
					case "/robots.txt":
						entry.Observations, entry.Recommendations, entry.sitemapURLs = analyzeRobotsTxt(body)
					case "/sitemap.xml":
						entry.Observations, entry.Recommendations = analyzeSitemapXml(body)
					case "/.well-known/security.txt":
						entry.Observations, entry.Recommendations = analyzeSecurityTxt(body)
					}
				}
			}

			mu.Lock()
			result.files = append(result.files, entry)
			mu.Unlock()
		}(fc)
	}

	wg.Wait()

	// Cross-reference: check that robots.txt Sitemap directive matches found sitemap.xml
	expectedSitemapURL := schema.String() + "://" + target.Hostname() + "/sitemap.xml"
	robotsIdx, sitemapIdx := -1, -1
	for i, f := range result.files {
		switch f.Path {
		case "/robots.txt":
			robotsIdx = i
		case "/sitemap.xml":
			sitemapIdx = i
		}
	}
	if robotsIdx >= 0 && sitemapIdx >= 0 && result.files[robotsIdx].Found && result.files[sitemapIdx].Found {
		referenced := false
		for _, u := range result.files[robotsIdx].sitemapURLs {
			if u == expectedSitemapURL {
				referenced = true
				break
			}
		}
		if !referenced && len(result.files[robotsIdx].sitemapURLs) > 0 {
			result.files[robotsIdx].Recommendations = append(
				result.files[robotsIdx].Recommendations,
				"Sitemap directive doesn't reference the found sitemap.xml at "+expectedSitemapURL,
			)
		}
	}

	status.SpinningComplete(fmt.Sprintf("Scan of %s well-known files complete.", schema.String()))

	slog.Debug("knownFilesScan: Scan completed")

	return result
}
