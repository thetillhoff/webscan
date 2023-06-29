package webscan

import (
	protocolScan "github.com/thetillhoff/webscan/pkg/protocolScan"
)

func (engine Engine) ScanHttpHeaders() (Engine, error) {

	engine.httpHeaderRecommendations = protocolScan.GenerateHeaderRecommendations(engine.response)

	return engine, nil
}
