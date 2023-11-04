package webscan

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
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
	httpContentHtmlSizekB       float64
	httpContentInlineStyleSize  int
	httpContentInlineScriptSize int
	httpContentScriptSizes      map[string]float64
	httpContentStylesheetSizes  map[string]float64

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
