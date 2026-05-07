package httpProtocolScan

import (
	"log/slog"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	isAvailableViaPort80  bool
	isAvailableViaPort443 bool
	timeout               time.Duration
}

type ConfigOption = types.Option[scanConfig]

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

// WithTimeout sets the per-request timeout for HTTP version checks
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(sc *scanConfig) {
		sc.timeout = timeout
	}
}

func Scan(target types.Target, status *status.Status, options ...ConfigOption) (Result, error) {
	var (
		result = Result{}
	)

	config := &scanConfig{
		timeout: 5 * time.Second,
	}
	types.ApplyOptions(config, options)

	slog.Debug("httpProtocolScan: Scan started")

	status.SpinningUpdate("Scanning http protocols...")

	// Check redirects for HTTP and HTTPS in parallel
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		target80 := target
		target80.OverrideSchema(types.HTTP)
		statusCode, redirectLocation, err := CheckHTTPRedirects(target80, config.timeout)
		result.httpStatusCode = statusCode
		result.httpRedirectLocation = redirectLocation
		result.isAvailableViaHttp = err == nil
	}()

	go func() {
		defer wg.Done()
		target443 := target
		target443.OverrideSchema(types.HTTPS)
		statusCode, redirectLocation, err := CheckHTTPRedirects(target443, config.timeout)
		result.httpsStatusCode = statusCode
		result.httpsRedirectLocation = redirectLocation
		result.isAvailableViaHttps = err == nil
	}()

	wg.Wait()

	// Check HTTP versions for port 80 and port 443 in parallel
	var wg2 sync.WaitGroup
	var httpVersions, httpsVersions []string
	var httpErr, httpsErr error

	if config.isAvailableViaPort80 && result.isAvailableViaHttp {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			target80 := target
			target80.OverrideSchema(types.HTTP)
			httpVersions, httpErr = CheckHTTPVersions(target80, config.timeout)
		}()
	}

	if config.isAvailableViaPort443 && result.isAvailableViaHttps {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			target443 := target
			target443.OverrideSchema(types.HTTPS)
			httpsVersions, httpsErr = CheckHTTPVersions(target443, config.timeout)
		}()
	}

	wg2.Wait()

	if httpErr != nil {
		return result, httpErr
	}
	if httpsErr != nil {
		return result, httpsErr
	}

	result.httpVersions = httpVersions
	result.httpsVersions = httpsVersions

	result.recommendations = httpProtocolRecommendations(target, result, config.isAvailableViaPort80, config.isAvailableViaPort443)

	status.SpinningComplete("Scan of http protocols completed.")

	slog.Debug("httpProtocolScan: Scan completed")

	return result, nil
}
