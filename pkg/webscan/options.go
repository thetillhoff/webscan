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

// IsEmpty reports whether no scan group is enabled.
func (o ScanOptions) IsEmpty() bool {
	return o == ScanOptions{}
}

// EnableWeb enables the HTTP-facing scan groups (protocol, header, content, known files).
func (o *ScanOptions) EnableWeb() {
	o.HTTPProtocol = true
	o.HTTPHeader = true
	o.HTMLContent = true
	o.KnownFiles = true
}
