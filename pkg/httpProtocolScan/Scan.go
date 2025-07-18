package httpProtocolScan

import (
	"log/slog"

	"github.com/thetillhoff/webscan/v3/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	client                cachedHttpGetClient.Client
	isAvailableViaPort80  bool
	isAvailableViaPort443 bool
}

// ConfigOption represents a configuration option for HTTP protocol scanning
type ConfigOption func(*scanConfig)

// WithClient sets the client
func WithClient(client cachedHttpGetClient.Client) ConfigOption {
	return func(sc *scanConfig) {
		sc.client = client
	}
}

// WithIsAvailableViaPort80 sets the isAvailableViaPort80
func WithIsAvailableViaPort80(isAvailableViaPort80 bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.isAvailableViaPort80 = isAvailableViaPort80
	}
}

// WithIsAvailableViaPort443 sets the isAvailableViaPort443
func WithIsAvailableViaPort443(isAvailableViaPort443 bool) ConfigOption {
	return func(sc *scanConfig) {
		sc.isAvailableViaPort443 = isAvailableViaPort443
	}
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	var (
		err error

		result = Result{}
	)

	// Apply configuration options
	config := &scanConfig{}
	for _, option := range options {
		option(config)
	}

	slog.Debug("httpProtocolScan: Scan started")

	status.SpinningUpdate("Scanning http protocols...")

	// Scan HTTP for redirects
	target.OverrideSchema(types.HTTP)
	result.httpStatusCode, result.httpRedirectLocation, err = CheckHttpRedirects(target, config.client)
	result.isAvailableViaHttp = err == nil

	// Scan HTTPS for redirects
	target.OverrideSchema(types.HTTPS)
	result.httpsStatusCode, result.httpsRedirectLocation, err = CheckHttpRedirects(target, config.client)
	result.isAvailableViaHttps = err == nil

	// TODO check redirect from http zone apex to https www. prefix
	// TODO check redirect from https zone apex to https www. prefix
	// TODO check redirect from http www. prefix to https www. prefix
	// TODO check redirects to omit the port (it's unneeded if protocol is set and it's the default 80 or 443)

	// TODO follow redirects if desired -> Probably not here, but in Scan().
	// TODO Only check http versions when there is no redirect happening

	// TODO check that redirectlocations either end with a `/` or with a filename (e.g. `index.html`).

	if config.isAvailableViaPort80 && result.isAvailableViaHttp {
		target.OverrideSchema(types.HTTP)

		// Scan http versions
		result.httpVersions, err = CheckHttpVersions(target)
		if err != nil {
			return result, err
		}
	}

	if config.isAvailableViaPort443 && result.isAvailableViaHttps {
		target.OverrideSchema(types.HTTPS)

		// Scan https versions
		result.httpsVersions, err = CheckHttpVersions(target)
		if err != nil {
			return result, err
		}
	}

	result.recommendations = GetHttpProtocolRecommendationsForResult(target, result, config.isAvailableViaPort80, config.isAvailableViaPort443)

	status.SpinningComplete("Scan of http protocols completed.")

	slog.Debug("httpProtocolScan: Scan completed")

	return result, nil
}
