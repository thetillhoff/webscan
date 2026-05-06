package httpProtocolScan

import (
	"log/slog"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

func CheckHTTPVersions(target types.Target, timeout time.Duration) ([]string, error) {
	slog.Debug("httpProtocolScan: Checking available http versions started", "url", target.UrlString())

	type versionResult struct {
		version string
		err     error
	}

	var (
		wg      sync.WaitGroup
		results [3]versionResult
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		results[0].version, results[0].err = checkHTTP1(target, timeout)
	}()

	go func() {
		defer wg.Done()
		results[1].version, results[1].err = checkHTTP2(target, timeout)
	}()

	go func() {
		defer wg.Done()
		results[2].version, results[2].err = checkHTTP3(target, timeout)
	}()

	wg.Wait()

	availableHttpVersions := []string{}
	for _, r := range results {
		if r.err == nil {
			availableHttpVersions = append(availableHttpVersions, r.version)
		}
	}

	slog.Debug("httpProtocolScan: Checking available http versions completed", "url", target.UrlString())

	return availableHttpVersions, nil
}
