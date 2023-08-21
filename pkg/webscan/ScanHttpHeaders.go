package webscan

import (
	"fmt"

	protocolScan "github.com/thetillhoff/webscan/pkg/protocolScan"
)

func (engine Engine) ScanHttpHeaders() (Engine, error) {

	fmt.Println("Scanning HTTP headers...")

	engine.httpHeaderRecommendations = protocolScan.GenerateHeaderRecommendations(engine.response)

	return engine, nil
}
