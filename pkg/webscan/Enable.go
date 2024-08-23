package webscan

func (engine *Engine) EnableAllScansIfNoneAreExplicitlySet() {

	if !(engine.advancedDnsScan ||
		engine.ipScan ||
		engine.advancedPortScan ||
		engine.tlsScan ||
		engine.httpProtocolScan ||
		engine.httpHeaderScan ||
		engine.htmlContentScan ||
		engine.mailConfigScan ||
		engine.subDomainScan) { // If no Scans are enabled, enable all by default

		engine.EnableAllScans()
	}
}

func (engine *Engine) EnableAllScans() {
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

func (engine *Engine) EnableWebScans() {
	engine.httpProtocolScan = true
	engine.httpHeaderScan = true
	engine.htmlContentScan = true
}

func (engine *Engine) EnableDetailedDnsScan() {
	engine.advancedDnsScan = true
}

func (engine *Engine) EnableIpScan() {
	engine.ipScan = true
}

func (engine *Engine) EnableDetailedPortScan() {
	engine.advancedPortScan = true
}

func (engine *Engine) EnableTlsScan() {
	engine.tlsScan = true
}

func (engine *Engine) EnableHttpProtocolScan() {
	engine.httpProtocolScan = true
}

func (engine *Engine) EnableHttpHeaderScan() {
	engine.httpHeaderScan = true
}

func (engine *Engine) EnableHtmlContentScan() {
	engine.htmlContentScan = true
}

func (engine *Engine) EnableMailConfigScan() {
	engine.mailConfigScan = true
}

func (engine *Engine) EnableSubdomainScan() {
	engine.subDomainScan = true
}
