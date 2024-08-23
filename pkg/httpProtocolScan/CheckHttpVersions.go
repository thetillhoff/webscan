package httpProtocolScan

import "log/slog"

// Takes url
// and whether http and/or https should be checked
// and checks which HTTP versions the server can speak for each
// Return available versions on http, then https and finally a potential error
func CheckHttpVersions(schema string, url string) ([]string, error) {
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

	slog.Debug("httpProtocolScan: Getting http versions started")

	httpVersion1, err = checkHttp1(schema + url)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion1)
	}

	httpVersion2, err = checkHttp2(schema + url)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion2)
	}

	httpVersion3, err = checkHttp3(schema + url)
	if err == nil {
		availableHttpVersions = append(availableHttpVersions, httpVersion3)
	}

	slog.Debug("httpProtocolScan: Getting http versions completed")

	return availableHttpVersions, nil
}
