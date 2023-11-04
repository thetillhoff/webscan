package webscan

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
)

// inputUrl can be domain or IPv4 or IPv6
// dnsServer can be empty string
func (engine Engine) Scan(inputUrl string) (Engine, error) {
	var (
		err error

		client  *http.Client
		request *http.Request
	)

	engine.input = inputUrl

	if net.ParseIP(inputUrl) == nil { // If inputUrl is domain, scan dns and ips
		engine.inputType = Domain
		if engine.Verbose {
			fmt.Println("Input identified as Domain.")
		}

		if engine.Verbose {
			if engine.dnsServer != "" {
				fmt.Println("Using custom dns server:", engine.dnsServer)
			} else {
				fmt.Println("Using system dns server")
			}
		}

		if engine.DetailedDnsScan {
			engine, err = engine.ScanDnsDetailed(inputUrl)
			if err != nil {
				return engine, err
			}
		} else {
			engine, err = engine.ScanDnsSimple(inputUrl)
			if err != nil {
				return engine, err
			}
		}
	} else { // If inputUrl is IPaddress, don't scan dns and ips
		if dnsScan.IsIpv4(inputUrl) { // If inputUrl is ipv4 address
			engine.inputType = IPv4
			if engine.Verbose {
				fmt.Println("Input identified as IPv4 address.")
			}
			engine.dnsScanEngine.ARecords = append(engine.dnsScanEngine.ARecords, inputUrl)
		} else { // If inputUrl is ipv6 address
			engine.inputType = IPv6
			if engine.Verbose {
				fmt.Println("Input identified as IPv6 address.")
			}
			engine.dnsScanEngine.AAAARecords = append(engine.dnsScanEngine.AAAARecords, inputUrl)
		}
	}

	if engine.IpScan {
		engine, err = engine.ScanIps()
		if err != nil {
			return engine, err
		}
	}

	if engine.DetailedPortScan {
		engine, err = engine.ScanPortDetailed()
		if err != nil {
			return engine, err
		}
	} else {
		engine, err = engine.ScanPortSimple()
		if err != nil {
			return engine, err
		}
	}

	if engine.TlsScan {
		engine, err = engine.ScanTls(inputUrl)
		if err != nil {
			return engine, err
		}
	}

	if engine.HttpProtocolScan {
		engine, err = engine.ScanHttpProtocols(inputUrl)
		if err != nil {
			return engine, err
		}
	}

	if engine.isAvailableViaHttps && (engine.HttpHeaderScan || engine.HttpContentScan) {
		client = &http.Client{
			Timeout: 5 * time.Second, // TODO 5s might be a bit long?
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}, // Don't follow redirects // TODO Should we follow redirects or not?
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Ignore invalid tls certificates here (certificates are checked in another step, and might be interesting what's behind it anyway)
			},
		}
		request, err = http.NewRequest("GET", "https://"+inputUrl, nil) // Only for https pages.
		if err != nil {
			return engine, err
		}
		request.Header.Set("User-Agent", "Go-http-client/1.1") // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
		// TODO request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0") // Set "random" valid user agent to prevent bot-detection (as it happens f.e. at amazon.com)
		engine.response, err = client.Do(request)
		if err != nil {
			fmt.Println(err, reflect.TypeOf(err))
			return engine, err
		}
	}

	if engine.isAvailableViaHttps && engine.HttpHeaderScan {
		engine, err = engine.ScanHttpHeaders()
		if err != nil {
			return engine, err
		}
	}

	if engine.isAvailableViaHttps && engine.HttpContentScan {
		engine, err = engine.ScanHttpContent(inputUrl)
		if err != nil {
			return engine, err
		}
	}

	if engine.MailConfigScan {
		engine, err = engine.ScanMailConfig(inputUrl)
		if err != nil {
			return engine, err
		}
	}

	if engine.SubdomainScan {
		engine, err = engine.ScanSubdomains(inputUrl)
		if err != nil {
			return engine, err
		}
	}

	// TODO if followRedirects is true, CNAMEs should be followed (scan them, too)
	// TODO if followRedirects is true, http and https redirects should be followed (scan them, too)

	return engine, nil
}
