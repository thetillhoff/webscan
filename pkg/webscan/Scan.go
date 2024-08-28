package webscan

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/thetillhoff/webscan/v3/pkg/dnsScan"
	"github.com/thetillhoff/webscan/v3/pkg/htmlContentScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpHeaderScan"
	"github.com/thetillhoff/webscan/v3/pkg/httpProtocolScan"
	"github.com/thetillhoff/webscan/v3/pkg/ipScan"
	"github.com/thetillhoff/webscan/v3/pkg/portScan"
	"github.com/thetillhoff/webscan/v3/pkg/subDomainScan"
	"github.com/thetillhoff/webscan/v3/pkg/tlsScan"
)

func (engine *Engine) Scan(input string) error {
	var (
		err error

		httpResponse  *http.Response
		httpsResponse *http.Response

		httpResponseBody  []byte
		httpsResponseBody []byte
	)

	// TODO If tty supports color, use custom logger, else use structured logger with zerolog or slog

	// TODO move the following test lines into a test case for `status` package. Ensure the final output looks the same as expected.
	// engine.status.Update("working on a")
	// time.Sleep(time.Second)
	// engine.status.Update("working on b")
	// time.Sleep(time.Second)
	// engine.status.Update("working on c")
	// time.Sleep(time.Second)
	// engine.status.Complete("complete")

	// engine.status.SpinningUpdate("working on a")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningUpdate("working on b")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningUpdate("working on c")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningComplete("complete")

	// engine.status.SpinningXOfInit(3, "things completed")
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// time.Sleep(3 * time.Second)
	// engine.status.SpinningXOfUpdate()
	// engine.status.SpinningXOfComplete()

	// Debug

	// TODO merge into one .Debug call
	slog.Debug("advancedDnsScan: " + strconv.FormatBool(engine.advancedDnsScan))
	slog.Debug("ipScan: " + strconv.FormatBool(engine.ipScan))
	slog.Debug("advancedPortScan: " + strconv.FormatBool(engine.advancedPortScan))
	slog.Debug("tlsScan: " + strconv.FormatBool(engine.tlsScan))
	slog.Debug("httpProtocolScan: " + strconv.FormatBool(engine.httpProtocolScan))
	slog.Debug("httpHeaderScan: " + strconv.FormatBool(engine.httpHeaderScan))
	slog.Debug("htmlContentScan: " + strconv.FormatBool(engine.htmlContentScan))
	slog.Debug("mailConfigScan: " + strconv.FormatBool(engine.mailConfigScan))
	slog.Debug("subDomainScan: " + strconv.FormatBool(engine.subDomainScan))

	// Input

	slog.Debug("input: " + input)
	engine.target, err = NewTarget(input)
	if err != nil {
		return err
	}
	slog.Debug("target", "engine.target", engine.target)

	// Output header for --instant
	if engine.instant {
		fmt.Printf("# webscan results for %s\n", engine.target.rawTarget)
	}

	// DNS

	if engine.target.isDomain {

		slog.Info("hostname identified as domain")

		// Since hostname is domain, dns is relevant

		if engine.advancedDnsScan { // If advanced scan is selected
			engine.dnsScanResult, err = dnsScan.AdvancedScan(
				&engine.status,
				engine.resolver,
				engine.target.Hostname(),
				engine.opinionated,
				engine.followRedirects,
			)
			if err != nil {
				return err
			}
		} else { // If simple scan is selected
			engine.dnsScanResult, err = dnsScan.SimpleScan(
				engine.resolver,
				engine.target.Hostname(),
				engine.followRedirects,
			)
			if err != nil {
				return err
			}
		}
	}

	if engine.instant {
		slog.Debug("instant dns result")
		dnsScan.PrintResult(engine.dnsScanResult)
	}

	if engine.ipScan {
		engine.ipScanResult, err = ipScan.Scan(
			&engine.status,
			engine.dnsScanResult.ARecords,
			engine.dnsScanResult.AAAARecords,
		)
		if err != nil {
			return err
		}
	}

	if engine.instant {
		slog.Debug("instant ip result")
		ipScan.PrintResult(
			engine.ipScanResult,
			engine.dnsScanResult.ARecords,
			engine.dnsScanResult.AAAARecords,
		)
	}

	if engine.advancedPortScan { // If advanced scan is selected
		engine.portScanResult, err = portScan.AdvancedScan(
			&engine.status,
			engine.dnsScanResult.ARecords,
			engine.dnsScanResult.AAAARecords,
		)
		if err != nil {
			return err
		}
	} else { // If simple scan is selected
		engine.portScanResult, err = portScan.SimpleScan(
			&engine.status,
			engine.dnsScanResult.ARecords,
			engine.dnsScanResult.AAAARecords,
		)
		if err != nil {
			return err
		}
	}

	if engine.instant {
		slog.Debug("instant port result")
		portScan.PrintResult(engine.portScanResult)
	}

	if engine.portScanResult.IsAvailableViaHttps() {

		if engine.tlsScan || engine.subDomainScan {
			engine.tlsScanResult, err = tlsScan.Scan(
				&engine.status,
				input,
				engine.target.parsedUrl,
			)
			if err != nil {
				return err
			}
		}

		if engine.instant {
			slog.Debug("instant tls result")
			tlsScan.PrintResult(engine.tlsScanResult)
		}
	}

	if engine.portScanResult.IsAvailableViaHttp() || engine.portScanResult.IsAvailableViaHttps() { // Only scan http protocol if target is reachable ;)

		if engine.httpProtocolScan {
			engine.httpProtocolScanResult, err = httpProtocolScan.Scan(
				&engine.status,
				input,
				engine.portScanResult.IsAvailableViaHttp(),
				engine.portScanResult.IsAvailableViaHttps(),
			)
			if err != nil {
				return err
			}
		}

		if engine.instant {
			slog.Debug("instant httpProtocol result")
			httpProtocolScan.PrintResult(engine.httpProtocolScanResult)
		}

	}

	if engine.httpHeaderScan || engine.htmlContentScan {

		// Make http request for later analysis of response
		if engine.portScanResult.IsAvailableViaHttp() {
			httpResponse, err = engine.client.MakeRequest("GET", "http://"+input, nil)
			if err != nil {
				return err
			}

			httpResponseBody, err = io.ReadAll(httpResponse.Body) // Read the body
			if err != nil {
				return err
			}
			defer httpResponse.Body.Close()
		}

		// Make https request for later analysis of response
		if engine.portScanResult.IsAvailableViaHttps() {
			httpsResponse, err = engine.client.MakeRequest("GET", "https://"+input, nil)
			if err != nil {
				return err
			}

			httpsResponseBody, err = io.ReadAll(httpsResponse.Body) // Read the body
			if err != nil {
				return err
			}
			defer httpsResponse.Body.Close()
		}

	}

	if engine.httpHeaderScan {

		if engine.portScanResult.IsAvailableViaHttp() {
			engine.httpHeaderScanResult = httpHeaderScan.Scan(&engine.status, httpResponse, "http")

			if engine.instant {
				slog.Debug("instant httpHeader result for http")
				httpHeaderScan.PrintResult(engine.httpHeaderScanResult, "http")
			}
		}

		if engine.portScanResult.IsAvailableViaHttps() {
			engine.httpsHeaderScanResult = httpHeaderScan.Scan(&engine.status, httpsResponse, "https")

			if engine.instant {
				slog.Debug("instant httpHeader result for https")
				httpHeaderScan.PrintResult(engine.httpsHeaderScanResult, "https")
			}
		}
	}

	if engine.htmlContentScan {

		if engine.portScanResult.IsAvailableViaHttp() {
			engine.httpHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, input, httpResponse, httpResponseBody, engine.client, "http")
			if err != nil {
				return err
			}

			if engine.instant {
				slog.Debug("instant htmlContent result for http")
				htmlContentScan.PrintResult(engine.httpHtmlContentScanResult, "http")
			}
		}

		if engine.portScanResult.IsAvailableViaHttps() {
			engine.httpsHtmlContentScanResult, err = htmlContentScan.Scan(&engine.status, input, httpsResponse, httpsResponseBody, engine.client, "https")
			if err != nil {
				return err
			}

			if engine.instant {
				slog.Debug("instant htmlContent result for https")
				htmlContentScan.PrintResult(engine.httpsHtmlContentScanResult, "https")
			}
		}

	}

	// if engine.MailConfigScan {
	// 	engine, err = engine.ScanMailConfig(input)
	// 	if err != nil {
	// 		return engine, err
	// 	}
	// }

	if engine.subDomainScan {
		engine.subDomainScanResult = subDomainScan.Scan(&engine.status, input, engine.tlsScanResult.CertNames)
	}

	if engine.instant {
		slog.Debug("instant subDomain result")
		subDomainScan.PrintResult(engine.subDomainScanResult)
	}

	// TODO if followRedirects is true, CNAMEs should be followed (scan them, too)
	// TODO if followRedirects is true, http and https redirects should be followed (scan them, too)

	return err
}
