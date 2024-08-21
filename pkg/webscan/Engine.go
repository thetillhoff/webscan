package webscan

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
	"github.com/thetillhoff/webscan/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/pkg/httpClient"
	"github.com/thetillhoff/webscan/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/pkg/ipScan"
	"github.com/thetillhoff/webscan/pkg/portScan"
	"github.com/thetillhoff/webscan/pkg/status"
	"github.com/thetillhoff/webscan/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

// TODO add proper logger implementation, with info statements on -v, and debug statements on -vvv

// TODO remove "opinionated" flag: Either it's a valid recommendation or not.

type Engine2 struct {
	status status.Status

	target Target

	resolver *net.Resolver // Nil resolver (==nil) is the same as a zero resolver which is the default system resolver
	client   httpClient.Client

	// Global settings
	followRedirects bool
	instant         bool
	opinionated     bool

	// Enabled/Disabled scans
	advancedDnsScan  bool
	ipScan           bool
	advancedPortScan bool
	tlsScan          bool
	httpProtocolScan bool
	httpHeaderScan   bool
	htmlContentScan  bool
	mailConfigScan   bool
	subDomainScan    bool

	// Results
	dnsScanResult              dnsScan.Result
	ipScanResult               ipScan.Result
	portScanResult             portScan.Result
	tlsScanResult              tlsScan.Result
	httpProtocolScanResult     httpProtocolScan.Result
	httpHeaderScanResult       httpHeaderScan.Result
	httpsHeaderScanResult      httpHeaderScan.Result
	httpHtmlContentScanResult  htmlContentScan.Result
	httpsHtmlContentScanResult htmlContentScan.Result
	subDomainScanResult        subDomainScan.Result
	// mailConfigScanResults []string // TODO find better type
}

func NewEngine(
	quiet bool,
	noColor bool,
	dnsServer string,
	followRedirects bool,
	instant bool,
	advancedDnsScan bool,
	ipScan bool,
	advancedPortScan bool,
	tlsScan bool,
	httpProtocolScan bool,
	httpHeaderScan bool,
	htmlContentScan bool,
	mailConfigScan bool,
	subDomainScan bool,
	writeMutex *sync.Mutex,
) (Engine2, error) {

	var (
		engine   Engine2
		resolver *net.Resolver
		client   httpClient.Client
	)

	if dnsServer != "" {
		resolver = &net.Resolver{
			PreferGo:     false,
			StrictErrors: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 10 * time.Second} // 10s timeout by default // TODO make this variable
				return d.DialContext(ctx, network, net.JoinHostPort(dnsServer, "53"))
			},
		}
		slog.Info("Using custom dns server", "dnsServer", dnsServer) // TODO use INF
	} else {
		slog.Info("Using system dns server") // TODO use INF
	}

	client = httpClient.NewClient(
		5*time.Second,
		10,
		false,
		"Go-http-client/1.1",
	)

	engine = Engine2{
		status:           status.NewStatus(quiet, noColor, writeMutex),
		resolver:         resolver,
		client:           client,
		followRedirects:  followRedirects,
		instant:          instant,
		advancedDnsScan:  advancedDnsScan,
		ipScan:           ipScan,
		advancedPortScan: advancedPortScan,
		tlsScan:          tlsScan,
		httpProtocolScan: httpProtocolScan,
		httpHeaderScan:   httpHeaderScan,
		htmlContentScan:  htmlContentScan,
		mailConfigScan:   mailConfigScan,
		subDomainScan:    subDomainScan,
	}

	return engine, nil
}

// TODO remove below

// type Engine struct {
// 	// Input
// 	DkimSelector string

// 	// Settings
// 	Opinionated     bool
// 	Verbose         bool
// 	FollowRedirects bool

// 	// Scan modes
// 	DetailedDnsScan  bool
// 	IpScan           bool
// 	advancedPortScan bool
// 	TlsScan          bool
// 	HttpProtocolScan bool
// 	HttpHeaderScan   bool
// 	htmlContentScan  bool
// 	MailConfigScan   bool
// 	SubdomainScan    bool

// 	// Internal variables
// 	input         string
// 	inputType     InputType
// 	dnsServer     string
// 	resolver      *net.Resolver
// 	httpResponse  *http.Response
// 	httpsResponse *http.Response
// 	response      *http.Response
// 	client        httpClient.Client

// 	// Results
// 	// dnsScanEngine dnsScan.Engine

// 	ipOwners          []string
// 	ipIsBlacklistedAt map[string][]string

// 	openPorts               []uint16
// 	openPortInconsistencies []string

// 	isAvailableViaHttp    bool
// 	isAvailableViaHttps   bool
// 	httpStatusCode        int
// 	httpRedirectLocation  string
// 	httpsStatusCode       int
// 	httpsRedirectLocation string
// 	httpVersions          []string
// 	httpsVersions         []string

// 	httpHeaderRecommendations      []string
// 	httpCookieRecommendations      map[string][]string
// 	httpOtherCookieRecommendations []string

// 	tlsResult  error
// 	tlsCiphers []tlsScan.TlsCipher

// 	httpContentRecommendations  []string
// 	httpContentHtmlSize         int
// 	httpContentInlineStyleSize  int
// 	httpContentInlineScriptSize int
// 	httpContentScriptSizes      map[string]int
// 	httpContentStylesheetSizes  map[string]int

// 	subdomains []string

// 	mailConfigRecommendations []string
// }

// func DefaultEngine(inputUrl string) Engine {
// 	return Engine{
// 		Opinionated:     true,
// 		Verbose:         false,
// 		FollowRedirects: false,

// 		dnsScanEngine: dnsScan.DefaultEngine(),
// 		dnsServer:     "",
// 		resolver:      nil, // Nil resolver is the same as a zero resolver which is the default system resolver

// 		ipIsBlacklistedAt: map[string][]string{},

// 		DetailedDnsScan:  false,
// 		IpScan:           false,
// 		advancedPortScan: false,
// 		TlsScan:          false,
// 		HttpProtocolScan: false,
// 		HttpHeaderScan:   false,
// 		htmlContentScan:  false,
// 		MailConfigScan:   false,
// 		SubdomainScan:    false,

// 		DkimSelector: "",
// 	}
// }

// func EngineWithCustomDns(inputUrl string, dnsServer string) Engine {
// 	engine := DefaultEngine(inputUrl)
// 	engine.dnsServer = dnsServer
// 	engine.resolver = &net.Resolver{
// 		PreferGo:     false,
// 		StrictErrors: true,
// 		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
// 			d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
// 			return d.DialContext(ctx, network, net.JoinHostPort(dnsServer, "53"))
// 		},
// 	}
// 	return engine
// }
