package httpHeaderScan

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/thetillhoff/webscan/v3/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

type scanConfig struct {
	schemaOverride types.Schema
	client         cachedHttpGetClient.Client
}

// ConfigOption represents a configuration option for DNS scanning
type ConfigOption func(*scanConfig)

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
			httpHeaderRecommendations:      []string{},
			httpCookieRecommendations:      map[string][]string{},
			httpOtherCookieRecommendations: []string{},
		}

		err      error
		response *http.Response
	)

	slog.Debug("httpHeaderScan: Scan started")

	config := &scanConfig{}
	for _, option := range options {
		option(config)
	}

	target.OverrideSchema(config.schemaOverride)

	status.SpinningUpdate(fmt.Sprintf("Scanning %s headers...", target.Schema().String()))

	response, _, err = config.client.Get(target.UrlString())
	if err != nil {
		return result, err
	}

	result.httpHeaderRecommendations = append(result.httpHeaderRecommendations, GenerateHeaderRecommendations(response)...)

	result.httpCookieRecommendations, result.httpOtherCookieRecommendations = GenerateCookieRecommendations(response) // TODO append instead

	status.SpinningComplete(fmt.Sprintf("Scan of %s headers completed.", target.Schema().String()))

	slog.Debug("httpHeaderScan: Scan completed")

	return result, nil
}
