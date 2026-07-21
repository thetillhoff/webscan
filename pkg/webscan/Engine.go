package webscan

import (
	"context"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/thetillhoff/webscan/v5/pkg/cachedHttpGetClient"
	"github.com/thetillhoff/webscan/v5/pkg/dnsScan"
	"github.com/thetillhoff/webscan/v5/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/v5/pkg/ipScan"
	"github.com/thetillhoff/webscan/v5/pkg/portScan"
	"github.com/thetillhoff/webscan/v5/pkg/status"
	"github.com/thetillhoff/webscan/v5/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/v5/pkg/tlsScan"
	"github.com/thetillhoff/webscan/v5/pkg/types"
)

// TODO add proper logger implementation, with info statements on -v, and debug statements on -vvv

// TODO remove "opinionated" flag: Either it's a valid recommendation or not.

type Engine struct {
	status status.Status

	target types.Target

	resolver *net.Resolver // Nil resolver (==nil) is the same as a zero resolver which is the default system resolver
	client   cachedHttpGetClient.Client

	stdout    io.Writer
	dnsServer string

	// Global settings
	followRedirects bool
	timeout         time.Duration

	// Enabled/Disabled scans
	advancedDnsScan  bool
	ipScan           bool
	advancedPortScan bool
	tlsScan          bool
	httpProtocolScan bool
	httpHeaderScan   bool
	htmlContentScan  bool
	knownFilesScan   bool
	mailConfigScan   bool
	subDomainScan    bool

	// Results
	dnsScanResult          dnsScan.Result
	ipScanResult           ipScan.Result
	portScanResult         portScan.Result
	tlsScanResult          tlsScan.Result
	httpProtocolScanResult httpProtocolScan.Result
	subDomainScanResult    subDomainScan.Result
	// mailConfigScanResults []string // TODO find better type
}

func NewEngine(
	stdout io.Writer,
	statusOut io.Writer,
	noColor bool,
	dnsServer string,
	followRedirects bool,
	timeout time.Duration,
	opts ScanOptions,
	writeMutex *sync.Mutex,
) (Engine, error) {

	var (
		engine   Engine
		resolver *net.Resolver
		client   cachedHttpGetClient.Client
	)

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	if dnsServer != "" {
		resolver = &net.Resolver{
			PreferGo:     false,
			StrictErrors: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: timeout}
				return d.DialContext(ctx, network, net.JoinHostPort(dnsServer, "53"))
			},
		}
		slog.Info("webscan: Using custom DNS server", "dnsServer", dnsServer)
	} else {
		slog.Info("webscan: Using system DNS server")
	}

	client = cachedHttpGetClient.NewClient(
		timeout,
		10,
		false,
		"Go-http-client/1.1",
	)

	engine = Engine{
		stdout:           stdout,
		dnsServer:        dnsServer,
		status:           status.NewStatus(noColor, writeMutex, statusOut),
		resolver:         resolver,
		client:           client,
		followRedirects:  followRedirects,
		timeout:          timeout,
		advancedDnsScan:  opts.AdvancedDNS,
		ipScan:           opts.IP,
		advancedPortScan: opts.AdvancedPort,
		tlsScan:          opts.TLS,
		httpProtocolScan: opts.HTTPProtocol,
		httpHeaderScan:   opts.HTTPHeader,
		htmlContentScan:  opts.HTMLContent,
		knownFilesScan:   opts.KnownFiles,
		mailConfigScan:   opts.MailConfig,
		subDomainScan:    opts.Subdomain,
	}

	return engine, nil
}
