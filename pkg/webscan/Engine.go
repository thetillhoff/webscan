package webscan

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
	"github.com/thetillhoff/webscan/pkg/httpClient"
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

type Engine struct {
	// Input
	DkimSelector string

	// Settings
	Opinionated     bool
	Verbose         bool
	FollowRedirects bool

	// Scan modes
	DetailedDnsScan  bool
	IpScan           bool
	DetailedPortScan bool
	TlsScan          bool
	HttpAllScans     bool
	HttpProtocolScan bool
	HttpHeaderScan   bool
	HttpContentScan  bool
	MailConfigScan   bool
	SubdomainScan    bool

	// Internal variables
	input     string
	inputType InputType
	dnsServer string
	resolver  *net.Resolver
	response  *http.Response // internal use only
	client    httpClient.Client

	// Results
	dnsScanEngine dnsScan.Engine

	ipOwners          []string
	ipIsBlacklistedAt map[string][]string

	openPorts               []uint16
	openPortInconsistencies []string

	isAvailableViaHttp    bool
	isAvailableViaHttps   bool
	httpStatusCode        int
	httpRedirectLocation  string
	httpsStatusCode       int
	httpsRedirectLocation string
	httpVersions          []string
	httpsVersions         []string

	httpHeaderRecommendations      []string
	httpCookieRecommendations      map[string][]string
	httpOtherCookieRecommendations []string

	tlsResult  error
	tlsCiphers []tlsScan.TlsCipher

	httpContentRecommendations  []string
	httpContentHtmlSize         int
	httpContentInlineStyleSize  int
	httpContentInlineScriptSize int
	httpContentScriptSizes      map[string]int
	httpContentStylesheetSizes  map[string]int

	subdomains []string

	mailConfigRecommendations []string
}

func DefaultEngine(inputUrl string) Engine {
	return Engine{
		Opinionated:     true,
		Verbose:         false,
		FollowRedirects: false,

		dnsScanEngine: dnsScan.DefaultEngine(),
		dnsServer:     "",
		resolver:      nil, // Nil resolver is the same as a zero resolver which is the default system resolver

		ipIsBlacklistedAt: map[string][]string{},

		DetailedDnsScan:  false,
		IpScan:           false,
		DetailedPortScan: false,
		TlsScan:          false,
		HttpProtocolScan: false,
		HttpHeaderScan:   false,
		HttpContentScan:  false,
		MailConfigScan:   false,
		SubdomainScan:    false,

		DkimSelector: "",
	}
}

func EngineWithCustomDns(inputUrl string, dnsServer string) Engine {
	engine := DefaultEngine(inputUrl)
	engine.dnsServer = dnsServer
	engine.resolver = &net.Resolver{
		PreferGo:     false,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
			return d.DialContext(ctx, network, net.JoinHostPort(dnsServer, "53"))
		},
	}
	return engine
}

func (engine *Engine) EnableDetailedDnsScan() {
	engine.DetailedDnsScan = true
}

func (engine *Engine) EnableIpScan() {
	engine.IpScan = true
}

func (engine *Engine) EnableDetailedPortScan() {
	engine.DetailedPortScan = true
}

func (engine *Engine) EnableTlsScan() {
	engine.TlsScan = true
}

func (engine *Engine) EnableHttpProtocolScan() {
	engine.HttpProtocolScan = true
}

func (engine *Engine) EnableHttpHeaderScan() {
	engine.HttpHeaderScan = true
}

func (engine *Engine) EnableHttpContentScan() {
	engine.HttpContentScan = true
}

func (engine *Engine) EnableAllHttpScans() {
	engine.EnableHttpProtocolScan()
	engine.EnableHttpHeaderScan()
	engine.EnableHttpContentScan()
}

func (engine *Engine) EnableMailConfigScan() {
	engine.MailConfigScan = true
}

func (engine *Engine) EnableSubdomainScan() {
	engine.SubdomainScan = true
}

func (engine *Engine) EnableAllScans() {
	engine.EnableDetailedDnsScan()
	engine.EnableIpScan()
	engine.EnableDetailedPortScan()
	engine.EnableTlsScan()
	engine.EnableHttpProtocolScan()
	engine.EnableHttpHeaderScan()
	engine.EnableHttpContentScan()
	engine.EnableMailConfigScan()
	engine.EnableSubdomainScan()
}

func (engine *Engine) EnableAllScansIfNoneAreExplicitlySet() {

	if !(engine.DetailedDnsScan ||
		engine.IpScan ||
		engine.DetailedPortScan ||
		engine.TlsScan ||
		engine.HttpProtocolScan ||
		engine.HttpHeaderScan ||
		engine.HttpContentScan ||
		engine.MailConfigScan ||
		engine.SubdomainScan) { // If no Scans are enabled, enable all by default

		engine.EnableAllScans()
	}
}
