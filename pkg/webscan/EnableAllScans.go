package webscan

func (engine *Engine2) EnableAllScans() {
	engine.advancedDnsScan = true
	engine.ipScan = true
	engine.advancedPortScan = true
	engine.tlsScan = true
	engine.httpProtocolScan = true
	engine.httpHeaderScan = true
	engine.htmlContentScan = true
	engine.mailConfigScan = true
	engine.subDomainScan = true
}
