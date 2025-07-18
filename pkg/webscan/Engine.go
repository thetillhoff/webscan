package webscan

import (
	"context"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v3/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v3/pkg/dnsScan"
	"github.com/thetillhoff/webscan/v3/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/v3/pkg/ipScan"
	"github.com/thetillhoff/webscan/v3/pkg/portScan"
	"github.com/thetillhoff/webscan/v3/pkg/status"
	"github.com/thetillhoff/webscan/v3/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/v3/pkg/tlsScan"
	"github.com/thetillhoff/webscan/v3/pkg/types"
)

// TODO add proper logger implementation, with info statements on -v, and debug statements on -vvv

// TODO remove "opinionated" flag: Either it's a valid recommendation or not.

type Engine struct {
	status status.Status

	target types.Target

	resolver *net.Resolver // Nil resolver (==nil) is the same as a zero resolver which is the default system resolver
	client   cachedHttpGetClient.Client

	stdout io.Writer

	// Global settings
	followRedirects bool

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
	stdout io.Writer,
	noColor bool,
	dnsServer string,
	followRedirects bool,
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
		client   cachedHttpGetClient.Client
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

	client = cachedHttpGetClient.NewClient(
		5*time.Second,
		10,
		false,
		"Go-http-client/1.1",
	)

	engine = Engine{
		stdout:           stdout,
		status:           status.NewStatus(noColor, writeMutex, stdout),
		resolver:         resolver,
		client:           client,
		followRedirects:  followRedirects,
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
