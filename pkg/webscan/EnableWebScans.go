package webscan

func (engine *Engine2) EnableWebScans() {
	engine.httpProtocolScan = true
	engine.httpHeaderScan = true
	engine.htmlContentScan = true
}
