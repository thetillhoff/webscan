package httpProtocolScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// Takes url
// and whether http and/or https should be checked
// and checks which HTTP versions the server can speak for each
// Return available versions on http, then https and finally a potential error
func CheckHttpVersions(target types.Target) ([]string, error) {
	var (
		err                   error
		availableHttpVersions = []string{}

		httpVersion1 string
		httpVersion2 string
		httpVersion3 string
	)

	// HTTP versions:
	// 0.9 -> obsolete
	// 1.0 -> obsolete
	// 1.1
	// 2
	// 3 QUIC

	slog.Debug("httpProtocolScan: Checking available http versions started", "url", target.UrlString())

	httpVersion1, err = checkHttp1(target)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion1)
	}

	httpVersion2, err = checkHttp2(target)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion2)
	}

	httpVersion3, err = checkHttp3(target)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion3)
	}

	slog.Debug("httpProtocolScan: Checking available http versions completed", "url", target.UrlString())

	return availableHttpVersions, nil
}
