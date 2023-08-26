package webscan

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/thetillhoff/webscan/pkg/dnsScan"
)

// func (engine Engine) Scan(followRedirects bool) (Engine, error) {
func (engine Engine) Scan(inputUrl string) (Engine, error) {
	var (
		err error

		client  *http.Client
		request *http.Request
	)

	if engine.Verbose {
		fmt.Println("Engine used:")
		engine.PrintEngine()
		fmt.Println()
	}

	netIP := net.ParseIP(inputUrl)
	if netIP == nil { // If inputUrl is IPaddress -> Only scan dns and ips if input is a domain, not an ip address
		if engine.Verbose {
			fmt.Println("Input identified as Domain.")
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
	} else {
		if dnsScan.IsIpv4(inputUrl) { // If inputUrl is ipv4 address
			engine.DnsScanEngine.ARecords = append(engine.DnsScanEngine.ARecords, inputUrl)
			if engine.Verbose {
				fmt.Println("Input identified as IPv4 address.")
			}
		} else { // If inputUrl is ipv6 address
			engine.DnsScanEngine.AAAARecords = append(engine.DnsScanEngine.AAAARecords, inputUrl)
			if engine.Verbose {
				fmt.Println("Input identified as IPv6 address.")
			}
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

	if engine.HttpHeaderScan || engine.HttpContentScan {
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
			return engine, err
		}
	}

	if engine.HttpHeaderScan {
		engine, err = engine.ScanHttpHeaders()
		if err != nil {
			return engine, err
		}
	}

	if engine.HttpContentScan {
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
