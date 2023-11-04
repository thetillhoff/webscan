package webscan

import (
	"fmt"

	httpHeaderScan "github.com/thetillhoff/webscan/pkg/httpHeaderScan"
)

func (engine Engine) ScanHttpHeaders() (Engine, error) {

	fmt.Println("Scanning HTTP headers...")

	engine.httpHeaderRecommendations = httpHeaderScan.GenerateHeaderRecommendations(engine.response)

	engine.httpCookieRecommendations, engine.httpOtherCookieRecommendations = httpHeaderScan.GenerateCookieRecommendations(engine.response)

	return engine, nil
}
