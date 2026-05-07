package webscan

func (engine *Engine) EnableAllScansIfNoneAreExplicitlySet() {

	if !engine.advancedDnsScan &&
		!engine.ipScan &&
		!engine.advancedPortScan &&
		!engine.tlsScan &&
		!engine.httpProtocolScan &&
		!engine.httpHeaderScan &&
		!engine.htmlContentScan &&
		!engine.knownFilesScan &&
		!engine.mailConfigScan &&
		!engine.subDomainScan { // If no Scans are enabled, enable all by default

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
	engine.knownFilesScan = true
	engine.mailConfigScan = true
	engine.subDomainScan = true
}

func (engine *Engine) EnableWebScans() {
	engine.httpProtocolScan = true
	engine.httpHeaderScan = true
	engine.htmlContentScan = true
	engine.knownFilesScan = true
}

func (engine *Engine) EnableDetailedDnsScan() {
	engine.advancedDnsScan = true
}

func (engine *Engine) EnableIPScan() {
	engine.ipScan = true
}

func (engine *Engine) EnableDetailedPortScan() {
	engine.advancedPortScan = true
}

func (engine *Engine) EnableTLSScan() {
	engine.tlsScan = true
}

func (engine *Engine) EnableHTTPProtocolScan() {
	engine.httpProtocolScan = true
}

func (engine *Engine) EnableHTTPHeaderScan() {
	engine.httpHeaderScan = true
}

func (engine *Engine) EnableHTMLContentScan() {
	engine.htmlContentScan = true
}

func (engine *Engine) EnableMailConfigScan() {
	engine.mailConfigScan = true
}

func (engine *Engine) EnableKnownFilesScan() {
	engine.knownFilesScan = true
}

func (engine *Engine) EnableSubdomainScan() {
	engine.subDomainScan = true
}
