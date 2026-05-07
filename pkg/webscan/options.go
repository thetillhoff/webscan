package webscan

// ScanOptions controls which scan groups are enabled.
type ScanOptions struct {
	AdvancedDNS  bool
	IP           bool
	AdvancedPort bool
	TLS          bool
	HTTPProtocol bool
	HTTPHeader   bool
	HTMLContent  bool
	KnownFiles   bool
	MailConfig   bool
	Subdomain    bool
}

func NewScanOptions(
	advancedDNS bool,
	ip bool,
	advancedPort bool,
	tls bool,
	httpProtocol bool,
	httpHeader bool,
	htmlContent bool,
	knownFiles bool,
	mailConfig bool,
	subdomain bool,
) ScanOptions {
	return ScanOptions{
		AdvancedDNS:  advancedDNS,
		IP:           ip,
		AdvancedPort: advancedPort,
		TLS:          tls,
		HTTPProtocol: httpProtocol,
		HTTPHeader:   httpHeader,
		HTMLContent:  htmlContent,
		KnownFiles:   knownFiles,
		MailConfig:   mailConfig,
		Subdomain:    subdomain,
	}
}

func AllScanOptions() ScanOptions {
	return NewScanOptions(true, true, true, true, true, true, true, true, true, true)
}

func WebScanOptions() ScanOptions {
	return NewScanOptions(true, true, true, true, true, true, true, true, false, false)
}

func (o ScanOptions) Apply(engine *Engine) {
	if o.AdvancedDNS {
		engine.EnableDetailedDnsScan()
	}
	if o.IP {
		engine.EnableIPScan()
	}
	if o.AdvancedPort {
		engine.EnableDetailedPortScan()
	}
	if o.TLS {
		engine.EnableTLSScan()
	}
	if o.HTTPProtocol {
		engine.EnableHTTPProtocolScan()
	}
	if o.HTTPHeader {
		engine.EnableHTTPHeaderScan()
	}
	if o.HTMLContent {
		engine.EnableHTMLContentScan()
	}
	if o.KnownFiles {
		engine.EnableKnownFilesScan()
	}
	if o.MailConfig {
		engine.EnableMailConfigScan()
	}
	if o.Subdomain {
		engine.EnableSubdomainScan()
	}
}
