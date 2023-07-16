package webscan

import (
	"net"
	"net/http"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
	"github.com/thetillhoff/webscan/pkg/tlsScan"
)

type Engine struct {
	// Input
	url string

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

	// Results
	dnsScanResult               dnsScan.Engine
	ipScanResult                []string
	ipScanOwners                []string
	portScanOpenPorts           []uint16
	portScanInconsistencies     []string
	isAvailableViaHttp          bool
	isAvailableViaHttps         bool
	httpStatusCode              int
	httpRedirectLocation        string
	httpsStatusCode             int
	httpsRedirectLocation       string
	tlsResult                   error
	tlsCiphers                  []tlsScan.TlsCipher
	httpVersions                []string
	httpsVersions               []string
	subdomains                  []string
	response                    *http.Response // internal use only
	httpHeaderRecommendations   []string
	httpContentRecommendations  []string
	httpContentHtmlSize         int
	httpContentInlineStyleSize  int
	httpContentInlineScriptSize int
	httpContentScriptSizes      map[string]float32
	httpContentStylesheetSizes  map[string]float32
	mailConfigRecommendations   []string

	DnsScanEngine dnsScan.Engine
	DkimSelector  string
}

func DefaultEngine(inputUrl string, dnsServer string) Engine {
	var (
		dnsScanEngine dnsScan.Engine
	)

	if dnsServer != "" {
		dnsScanEngine = dnsScan.EngineWithCustomDns(dnsServer)
	} else {
		dnsScanEngine = dnsScan.DefaultEngine()
	}

	netIP := net.ParseIP(inputUrl) // Try to parse inputUrl as IPaddress
	if netIP != nil {              // If inputUrl is IPaddress
		if dnsScan.IsIpv4(inputUrl) { // If it's an ipv4 address
			dnsScanEngine.ARecords = append(dnsScanEngine.ARecords, inputUrl) // Add as ipv4 address
		} else { // If it's an ipv6 address
			dnsScanEngine.AAAARecords = append(dnsScanEngine.AAAARecords, inputUrl) // Add as ipv6 address
		}
	}

	return Engine{
		url: inputUrl,

		Opinionated:     true,
		Verbose:         false,
		FollowRedirects: false,

		DetailedDnsScan:  false,
		IpScan:           false,
		DetailedPortScan: false,
		TlsScan:          false,
		HttpProtocolScan: false,
		HttpHeaderScan:   false,
		HttpContentScan:  false,
		MailConfigScan:   false,
		SubdomainScan:    false,

		DnsScanEngine: dnsScanEngine,
		DkimSelector:  "",
	}
}
