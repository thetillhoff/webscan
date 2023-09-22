package webscan

import (
	"net/http"

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
	inputType InputType

	// Results
	dnsScanEngine               dnsScan.Engine
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
	protocolRecommendations     []string
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
}

func DefaultEngine(inputUrl string, dnsServer string) Engine {
	return Engine{
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

		DkimSelector: "",
	}
}
