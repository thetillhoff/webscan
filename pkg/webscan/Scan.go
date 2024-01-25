package webscan

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
)

// input can be domain or IPv4 or IPv6
// dnsServer can be empty string
func (engine Engine) Scan(input string) (Engine, error) {
	var (
		err error

		hostname string
		client   *http.Client
		request  *http.Request
	)

	engine.input = input

	if net.ParseIP(input) == nil { // If input is domain, scan dns and ips
		if strings.Contains(input, "/") {
			hostname = strings.SplitN(input, "/", 2)[0]
		} else {
			hostname = input
		}

		if hostname == input { // If there is no path set in the input
			engine.inputType = DomainWithPath
			if engine.Verbose {
				fmt.Println("Input identified as Domain with path.")
			}
		} else { // If there is a path set in the input
			engine.inputType = Domain
			if engine.Verbose {
				fmt.Println("Input identified as Domain.")
			}
		}

		if engine.Verbose {
			if engine.dnsServer != "" {
				fmt.Println("Using custom dns server:", engine.dnsServer)
			} else {
				fmt.Println("Using system dns server")
			}
		}

		if engine.DetailedDnsScan {
			engine, err = engine.ScanDnsDetailed(hostname)
			if err != nil {
				return engine, err
			}
		} else {
			engine, err = engine.ScanDnsSimple(hostname)
			if err != nil {
				return engine, err
			}
		}
	} else { // If input is IPaddress, don't scan dns and ips
		if dnsScan.IsIpv4(input) { // If input is ipv4 address
			engine.inputType = IPv4
			if engine.Verbose {
				fmt.Println("Input identified as IPv4 address.")
			}
			engine.dnsScanEngine.ARecords = append(engine.dnsScanEngine.ARecords, input)
		} else { // If input is ipv6 address
			engine.inputType = IPv6
			if engine.Verbose {
				fmt.Println("Input identified as IPv6 address.")
			}
			engine.dnsScanEngine.AAAARecords = append(engine.dnsScanEngine.AAAARecords, input)
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
		engine, err = engine.ScanTls(input)
		if err != nil {
			return engine, err
		}
	}

	if engine.HttpProtocolScan {
		engine, err = engine.ScanHttpProtocols(input)
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
		request, err = http.NewRequest("GET", "https://"+input, nil) // Only for https pages.
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
		engine, err = engine.ScanHttpContent(input)
		if err != nil {
			return engine, err
		}
	}

	if engine.MailConfigScan {
		engine, err = engine.ScanMailConfig(input)
		if err != nil {
			return engine, err
		}
	}

	if engine.SubdomainScan {
		engine, err = engine.ScanSubdomains(input)
		if err != nil {
			return engine, err
		}
	}

	// TODO if followRedirects is true, CNAMEs should be followed (scan them, too)
	// TODO if followRedirects is true, http and https redirects should be followed (scan them, too)

	return engine, nil
}
