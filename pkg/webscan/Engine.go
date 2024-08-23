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

type Engine struct {
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
) (Engine, error) {

	var (
		engine   Engine
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

	engine = Engine{
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
