package httpHeaderScan

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/thetillhoff/webscan/v3/pkg/status"
)

func Scan(status *status.Status, response *http.Response, schemaName string) Result {
	var (
		result = Result{
			httpHeaderRecommendations:      []string{},
			httpCookieRecommendations:      map[string][]string{},
			httpOtherCookieRecommendations: []string{},
		}
	)

	slog.Debug("httpHeaderScan: Scan started")

	status.SpinningUpdate(fmt.Sprintf("Scanning %s headers...", schemaName))

	result.httpHeaderRecommendations = append(result.httpHeaderRecommendations, GenerateHeaderRecommendations(response)...)

	result.httpCookieRecommendations, result.httpOtherCookieRecommendations = GenerateCookieRecommendations(response) // TODO append instead

	status.SpinningComplete(fmt.Sprintf("Scan of %s headers completed.", schemaName))

	slog.Debug("httpHeaderScan: Scan completed")

	return result
}
