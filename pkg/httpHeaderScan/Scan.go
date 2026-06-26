package httpHeaderScan

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/thetillhoff/webscan/v5/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

type scanConfig struct {
	schemaOverride types.Schema
	client         cachedHttpGetClient.Client
}

type ConfigOption = types.Option[scanConfig]

// WithClient sets the client
func WithClient(client cachedHttpGetClient.Client) ConfigOption {
	return func(sc *scanConfig) {
		sc.client = client
	}
}

// WithSchemaOverride sets the schema override
func WithSchemaOverride(schema types.Schema) ConfigOption {
	return func(sc *scanConfig) {
		sc.schemaOverride = schema
	}
}

func Scan(status *status.Status, target types.Target, options ...ConfigOption) (Result, error) {
	var (
		result = Result{
			httpHeaderEntries:              []HeaderEntry{},
			httpCookieRecommendations:      map[string][]string{},
			httpOtherCookieRecommendations: []string{},
		}

		err      error
		response *http.Response
	)

	slog.Debug("httpHeaderScan: Scan started")

	config := &scanConfig{}
	types.ApplyOptions(config, options)

	target.OverrideSchema(config.schemaOverride)

	status.SpinningUpdate(fmt.Sprintf("Scanning %s headers...", target.Schema().String()))

	response, _, err = config.client.Get(target.UrlString())
	if err != nil {
		return result, err
	}

	result.httpHeaderEntries = GenerateHeaderRecommendations(response, config.schemaOverride)

	result.httpCookieRecommendations, result.httpOtherCookieRecommendations = GenerateCookieRecommendations(response) // TODO append instead

	status.SpinningComplete(fmt.Sprintf("Scan of %s headers completed.", target.Schema().String()))

	slog.Debug("httpHeaderScan: Scan completed")

	return result, nil
}
