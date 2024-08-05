package httpProtocolScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/pkg/status"
)

func Scan(status *status.Status, input string, isAvailableViaHttp bool, isAvailableViaHttps bool) (Result, error) {
	var (
		err error

		result = Result{}
	)

	slog.Debug("httpProtocolScan: Scan started")

	status.SpinningUpdate("Scanning http protocols...")

	// Scan HTTP / HTTPS for redirects
	result.httpStatusCode, result.httpRedirectLocation, result.httpsStatusCode, result.httpsRedirectLocation, err = CheckHttpRedirects(input, isAvailableViaHttp, isAvailableViaHttps)
	if err != nil {
		return result, err
	}

	// TODO check redirect from http zone apex to https www. prefix
	// TODO check redirect from https zone apex to https www. prefix
	// TODO check redirect from http www. prefix to https www. prefix
	// TODO check redirects to omit the port (it's unneeded if protocol is set and it's the default 80 or 443)

	// TODO follow redirects if desired -> Probably not here, but in Scan().
	// TODO Only check http versions when there is no redirect happening

	if isAvailableViaHttp && result.httpRedirectLocation != "" {
		// Scan http versions
		result.httpVersions, err = CheckHttpVersions("http://", input)
		if err != nil {
			return result, err
		}
	}

	if isAvailableViaHttps && result.httpsRedirectLocation != "" {
		// Scan https versions
		result.httpsVersions, err = CheckHttpVersions("https://", input)
		if err != nil {
			return result, err
		}
	}

	result.recommendations = GetHttpProtocolRecommendationsForResult(input, result, isAvailableViaHttp, isAvailableViaHttps)

	status.SpinningComplete("Scan of http protocols completed.")

	slog.Debug("httpProtocolScan: Scan completed")

	return result, nil
}
